package chain

import (
  "context"
  "golang.org/x/crypto/sha3"
  "encoding/json"
  "log"
  "math/big"
  "sync"
  "time"

  "github.com/ethereum/go-ethereum"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/ethclient"
  k "github.com/segmentio/kafka-go"

  "loyalty-points-system/internal/models"
)

type ChainCfg struct {
  Name           string `json:"name"`
  WSSURL         string `json:"wss_url"`
  TokenAddress   string `json:"token_address"`
  StakingAddress string `json:"staking_address"`
  Confirmations  int    `json:"confirmations"`
}
type Writer interface{ WriteMessages(ctx context.Context, msgs ...k.Message) error }

type Worker struct {
  cfg    ChainCfg
  client *ethclient.Client
  w      Writer

  transferTopic          common.Hash
  stakingTopicStaked     common.Hash
  stakingTopicWithdrawn  common.Hash
  stakingTopicReward     common.Hash
  stakingTopicEmergency  common.Hash
  stakingAddr            common.Address

  pending   map[uint64][]models.BalanceEvent
  pendingMu sync.Mutex
  head      uint64
}

func NewWorker(cfg ChainCfg, w Writer) *Worker {
  h := sha3.NewLegacyKeccak256(); h.Write([]byte("Transfer(address,address,uint256)"))
  transfer := common.BytesToHash(h.Sum(nil))
  hs := sha3.NewLegacyKeccak256(); hs.Write([]byte("Staked(address,uint256)")); stakingStaked := common.BytesToHash(hs.Sum(nil))
  hw := sha3.NewLegacyKeccak256(); hw.Write([]byte("Withdrawn(address,uint256)")); stakingWithdrawn := common.BytesToHash(hw.Sum(nil))
  hr := sha3.NewLegacyKeccak256(); hr.Write([]byte("RewardClaimed(address,uint256)")); stakingReward := common.BytesToHash(hr.Sum(nil))
  he := sha3.NewLegacyKeccak256(); he.Write([]byte("EmergencyWithdraw(address,uint256)")); stakingEmergency := common.BytesToHash(he.Sum(nil))

  return &Worker{
    cfg: cfg, w: w,
    transferTopic: transfer,
    stakingTopicStaked: stakingStaked,
    stakingTopicWithdrawn: stakingWithdrawn,
    stakingTopicReward: stakingReward,
    stakingTopicEmergency: stakingEmergency,
    pending: make(map[uint64][]models.BalanceEvent),
  }
}

func (wk *Worker) Start(ctx context.Context) error {
  var err error
  wk.client, err = ethclient.DialContext(ctx, wk.cfg.WSSURL)
  if err != nil { return err }
  log.Printf("🔌 [%s] connected", wk.cfg.Name)

  headers := make(chan *types.Header, 32)
  subHead, err := wk.client.SubscribeNewHead(ctx, headers)
  if err != nil { return err }

  var subLogs ethereum.Subscription
  logsCh := make(chan types.Log, 256)
  if wk.cfg.TokenAddress != "" && wk.cfg.TokenAddress != "0x0000000000000000000000000000000000000000" {
    addr := common.HexToAddress(wk.cfg.TokenAddress)
    q := ethereum.FilterQuery{Addresses: []common.Address{addr}, Topics: [][]common.Hash{{wk.transferTopic}}}
    subLogs, err = wk.client.SubscribeFilterLogs(ctx, q, logsCh)
    if err != nil { return err }
    log.Printf("👂 [%s] subscribed Transfer for %s", wk.cfg.Name, addr.Hex())
  } else {
    log.Printf("ℹ️  [%s] token_address not set; skip Transfer subscription", wk.cfg.Name)
  }
  if wk.cfg.StakingAddress != "" && wk.cfg.StakingAddress != "0x0000000000000000000000000000000000000000" {
    wk.stakingAddr = common.HexToAddress(wk.cfg.StakingAddress)
    q2 := ethereum.FilterQuery{
      Addresses: []common.Address{wk.stakingAddr},
      Topics: [][]common.Hash{{wk.stakingTopicStaked, wk.stakingTopicWithdrawn, wk.stakingTopicReward, wk.stakingTopicEmergency}},
    }
    sub2, err2 := wk.client.SubscribeFilterLogs(ctx, q2, logsCh)
    if err2 != nil { return err2 }
    log.Printf("👂 [%s] subscribed Staking logs for %s", wk.cfg.Name, wk.stakingAddr.Hex())
    go func(){ <-sub2.Err(); log.Printf("❌ [%s] staking logs sub err", wk.cfg.Name) }()
  } else {
    log.Printf("ℹ️  [%s] staking_address not set; skip staking subscription", wk.cfg.Name)
  }


  var logsErr <-chan error
  if subLogs != nil {
    logsErr = subLogs.Err()
  }
  conf := wk.cfg.Confirmations; if conf <= 0 { conf = 6 }

  go func() {
    for {
      select {
      case err := <-subHead.Err():
        log.Printf("❌ [%s] head sub err: %v", wk.cfg.Name, err); return
      case h := <-headers:
        if h != nil { wk.onHead(h.Number.Uint64(), conf) }
      case err := <-logsErr:
        if err!=nil { log.Printf("❌ [%s] logs sub err: %v", wk.cfg.Name, err); return }
      case lg := <-logsCh:
        wk.onLog(lg)
      case <-ctx.Done():
        return
      }
    }
  }()
  return nil
}

func (wk *Worker) onLog(lg types.Log) {
  // staking first
  if lg.Address == wk.stakingAddr && len(lg.Topics) > 0 {
    if len(lg.Topics) < 2 { return }
    user := common.BytesToAddress(lg.Topics[1].Bytes()[12:]).Hex()
    amt := new(big.Int).SetBytes(lg.Data).String()
    etype := ""
    switch lg.Topics[0] {
    case wk.stakingTopicStaked: etype = "transfer_in"
    case wk.stakingTopicWithdrawn, wk.stakingTopicEmergency: etype = "transfer_out"
    case wk.stakingTopicReward: etype = "transfer_in"
    default: return
    }
    evt := models.BalanceEvent{
      UserAddress: user, Amount: amt, EventType: etype,
      TxHash: lg.TxHash.Hex(), Chain: wk.cfg.Name, BlockNumber: int64(lg.BlockNumber),
      Confirmed: false, Timestamp: time.Now().Unix(),
    }
    wk.pendingMu.Lock(); wk.pending[lg.BlockNumber] = append(wk.pending[lg.BlockNumber], evt); wk.pendingMu.Unlock()
    log.Printf("🕒 [%s] staking pending %s user=%s amt=%s blk=%d", wk.cfg.Name, etype, user, amt, lg.BlockNumber)
    return
  }
  // ERC-20 Transfer
  if len(lg.Topics) == 0 || lg.Topics[0] != wk.transferTopic { return }
  if len(lg.Topics) < 3 { return }
  from := common.BytesToAddress(lg.Topics[1].Bytes()[12:])
  to   := common.BytesToAddress(lg.Topics[2].Bytes()[12:])
  amt  := new(big.Int).SetBytes(lg.Data).String()
  txHash := lg.TxHash.Hex()
  block  := lg.BlockNumber
  evts := []models.BalanceEvent{
    {UserAddress: from.Hex(), Amount: amt, EventType: "transfer_out", TxHash: txHash, Chain: wk.cfg.Name, BlockNumber: int64(block), Confirmed: false, Timestamp: time.Now().Unix()},
    {UserAddress: to.Hex(),   Amount: amt, EventType: "transfer_in",  TxHash: txHash, Chain: wk.cfg.Name, BlockNumber: int64(block), Confirmed: false, Timestamp: time.Now().Unix()},
  }
  wk.pendingMu.Lock(); wk.pending[block] = append(wk.pending[block], evts...); wk.pendingMu.Unlock()
  log.Printf("🕒 [%s] pending Transfer %s blk=%d from=%s to=%s amt=%s", wk.cfg.Name, txHash, block, from.Hex(), to.Hex(), amt)
}

func (wk *Worker) onHead(n uint64, conf int) {
  wk.head = n
  var out []models.BalanceEvent
  keep := make(map[uint64][]models.BalanceEvent)
  wk.pendingMu.Lock()
  for bn, list := range wk.pending {
    if n >= uint64(conf) && bn <= n-uint64(conf) { out = append(out, list...) } else { keep[bn] = list }
  }
  wk.pending = keep
  wk.pendingMu.Unlock()

  if len(out) == 0 { return }
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
  for _, evt := range out {
    evt.Confirmed = true
    b, _ := json.Marshal(evt)
    if err := wk.w.WriteMessages(ctx, k.Message{Value: b}); err != nil {
      log.Printf("write kafka err: %v", err)
    } else {
      log.Printf("✅ [%s] confirmed tx=%s block=%d", wk.cfg.Name, evt.TxHash, evt.BlockNumber)
    }
  }
}
