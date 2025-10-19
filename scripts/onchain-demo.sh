#!/usr/bin/env bash
set -euo pipefail

# ===========================
# ⚠️ 测试用默认配置（仅用于测试/演示）
# ===========================
ETH_RPC_URL="${ETH_RPC_URL:-https://eth-sepolia.g.alchemy.com/v2/FP6JOVxZoc4lDScODskcP}"
PRIVATE_KEY="${PRIVATE_KEY:-0x7d8f24a8e095e19eed0ec39dcd666e94ec20b6097a166f0ba21281e05d25adf9}"

TOKEN="${TOKEN:-0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B}"
STAKING="${STAKING:-0x3dBF997f45b7AF6EA274A7c4e9BaaBAC79989eed}"

# ✅ 默认测试接收地址（可用 TO_ADDR_DEFAULT 覆盖）
DEFAULT_TO="${TO_ADDR_DEFAULT:-0x1111111111111111111111111111111111111111}"

# 是否跟随 listener 日志（1=跟，0=不跟）
TAIL_LOGS="${TAIL_LOGS:-1}"

need() { command -v "$1" >/dev/null 2>&1 || { echo "缺少依赖：$1"; exit 1; }; }
need cast
command -v docker >/dev/null 2>&1 || true

ME="$(cast wallet address --private-key "$PRIVATE_KEY")"
echo "RPC=$ETH_RPC_URL"
echo "ACCOUNT=$ME"
echo "TOKEN=$TOKEN"
echo "STAKING=$STAKING"
echo "DEFAULT_TO=$DEFAULT_TO"
echo

toWei() { cast --to-wei "$1" ether; }
tx_link() { echo "https://sepolia.etherscan.io/tx/$1"; }
send() { cast send "$@" --rpc-url "$ETH_RPC_URL" --private-key "$PRIVATE_KEY" --async; }

start_logs() {
  if [[ "${TAIL_LOGS}" == "1" ]] && command -v docker >/dev/null 2>&1; then
    if docker ps --format '{{.Names}}' | grep -q '^loyalty-listener$'; then
      echo ">> 跟随 listener 日志（Ctrl+C 退出，仅影响日志）"
      docker logs --since 5s -f loyalty-listener &
      LOG_PID=$!; trap 'kill $LOG_PID 2>/dev/null || true' EXIT
    fi
  fi
}

do_transfer() {
  local to="${1:-$DEFAULT_TO}"
  local amount_eth="${2:-0.01}"
  local amount_wei; amount_wei="$(toWei "$amount_eth")"
  echo ">> ERC20 transfer -> $to, amount=$amount_eth ($amount_wei wei)"
  local tx; tx="$(send "$TOKEN" "transfer(address,uint256)" "$to" "$amount_wei")"
  echo "TX: $tx"; tx_link "$tx"
}

do_approve() {
  local amount_eth="${1:-0.01}"
  local amount_wei; amount_wei="$(toWei "$amount_eth")"
  echo ">> approve(STAKING, $amount_eth)"
  local tx; tx="$(send "$TOKEN" "approve(address,uint256)" "$STAKING" "$amount_wei")"
  echo "TX: $tx"; tx_link "$tx"
}

do_stake() {
  local amount_eth="${1:-0.01}"
  local amount_wei; amount_wei="$(toWei "$amount_eth")"
  echo ">> stake($amount_eth)"
  local tx; tx="$(send "$STAKING" "stake(uint256)" "$amount_wei")"
  echo "TX: $tx"; tx_link "$tx"
}

do_withdraw() {
  local amount_eth="${1:-0.005}"
  local amount_wei; amount_wei="$(toWei "$amount_eth")"
  echo ">> withdraw($amount_eth)"
  local tx; tx="$(send "$STAKING" "withdraw(uint256)" "$amount_wei")"
  echo "TX: $tx"; tx_link "$tx"
}

do_claim() {
  echo ">> claimRewards()"
  local tx; tx="$(send "$STAKING" "claimRewards()")"
  echo "TX: $tx"; tx_link "$tx"
}

do_demo() {
  local to="${1:-$DEFAULT_TO}"
  local amt="${2:-0.01}"
  local wd="${3:-0.005}"
  echo "== DEMO: approve -> transfer -> stake -> (等待确认) -> withdraw =="
  do_approve "$amt"
  do_transfer "$to" "$amt"
  do_stake "$amt"
  echo "提示：listener 等设定确认数（默认 6）。稍后可执行：$0 withdraw $wd"
}

usage() {
  cat <<USAGE
用法：
  $0 transfer [to_addr] [amount_eth]   # ERC20 转账（默认 to=$DEFAULT_TO, amount=0.01）
  $0 approve [amount_eth]              # 授权给 staking（默认 0.01）
  $0 stake [amount_eth]                # 质押（默认 0.01）
  $0 withdraw [amount_eth]             # 提取（默认 0.005）
  $0 claim                             # 领取奖励
  $0 demo [to_addr] [amount_eth] [withdraw_eth]

环境变量可覆盖：ETH_RPC_URL, PRIVATE_KEY, TOKEN, STAKING, TO_ADDR_DEFAULT, TAIL_LOGS
USAGE
  exit 1
}

start_logs
cmd="${1:-}"; shift || true
case "$cmd" in
  transfer)  do_transfer "${1:-}" "${2:-}";;
  approve)   do_approve "${1:-}";;
  stake)     do_stake "${1:-}";;
  withdraw)  do_withdraw "${1:-}";;
  claim)     do_claim;;
  demo)      do_demo "${1:-}" "${2:-}" "${3:-}";;
  *)         usage;;
esac

