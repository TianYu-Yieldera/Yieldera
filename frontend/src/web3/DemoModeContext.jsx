import { createContext, useContext, useState, useEffect } from 'react';
import { config } from '../config/env';

const API_BASE = config.api.baseUrl;

const DEFAULT_DEMO_DATA = {
  pfiTokens: 10000,
  points: 5000,
  stakedTokens: 0,
  stakingRewards: 0,
  defiDeposits: [],
  stablecoinMinted: 0,
  collateral: 0,
};

const DemoModeContext = createContext();

export function DemoModeProvider({ children }) {
  const [demoMode, setDemoMode] = useState(false);
  const [demoAddress, setDemoAddress] = useState(null);
  const [updateKey, setUpdateKey] = useState(0); // 用于强制刷新
  const [demoData, setDemoData] = useState(() => ({ ...DEFAULT_DEMO_DATA }));

  const safeParseNumber = (value, fallback = 0) => {
    if (value === null || value === undefined) return fallback;
    if (typeof value === 'number') {
      return Number.isFinite(value) ? value : fallback;
    }
    const parsed = parseFloat(value);
    return Number.isFinite(parsed) ? parsed : fallback;
  };

  const applySummaryToState = (summary) => {
    if (!summary) return;
    setDemoData(prev => {
      const next = {
        ...prev,
        points: safeParseNumber(summary.points, prev.points),
        pfiTokens: safeParseNumber(summary.token_balance, prev.pfiTokens),
        collateral: safeParseNumber(summary?.stablecoin?.collateral, prev.collateral),
        stablecoinMinted: safeParseNumber(summary?.stablecoin?.debt, prev.stablecoinMinted),
      };
      localStorage.setItem('demo_data', JSON.stringify(next));
      return next;
    });
    setUpdateKey(prev => prev + 1);
  };

  const refreshDemoFromServer = async (address) => {
    if (!address) return;
    try {
      const response = await fetch(`${API_BASE}/api/demo/summary?address=${address}`);
      if (!response.ok) {
        throw new Error('failed summary fetch');
      }
      const summary = await response.json();
      const isActive = Boolean(summary.is_demo) && Boolean(summary.demo_active);
      setDemoMode(isActive);
      setDemoAddress(summary.address || address);
      localStorage.setItem('demo_mode_active', String(isActive));
      localStorage.setItem('demo_wallet_address', summary.address || address);
      applySummaryToState(summary);
    } catch (err) {
      console.warn('Failed to refresh demo summary', err);
    }
  };

  // 计算实时收益的辅助函数
  const calculateEarnings = (deposit) => {
    if (!deposit.depositTime) return 0;
    const timeElapsed = (Date.now() - deposit.depositTime) / 1000; // 秒
    // 按秒计算收益：(amount * apy / 100) / (365 * 24 * 60 * 60) * timeElapsed
    return (deposit.amount * deposit.apy / 100) / (365 * 24 * 60 * 60) * timeElapsed;
  };

  // 获取带实时收益的存款数据
  const getDepositsWithEarnings = () => {
    return demoData.defiDeposits.map(deposit => ({
      ...deposit,
      earned: calculateEarnings(deposit)
    }));
  };

  useEffect(() => {
    const savedMode = localStorage.getItem('demo_mode_active') === 'true';
    const savedAddress = localStorage.getItem('demo_wallet_address');
    const savedData = localStorage.getItem('demo_data');

    if (savedData) {
      try {
        const parsedData = JSON.parse(savedData);
        setDemoData(prev => ({
          ...prev,
          ...DEFAULT_DEMO_DATA,
          ...parsedData,
        }));
      } catch (e) {
        console.warn('Failed to parse demo data');
      }
    }

    if (savedMode && savedAddress) {
      setDemoMode(true);
      setDemoAddress(savedAddress);
      refreshDemoFromServer(savedAddress);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const enableDemoMode = async () => {
    const randomAddress = '0x' + Array.from({ length: 40 }, () =>
      Math.floor(Math.random() * 16).toString(16)
    ).join('');

    try {
      const response = await fetch(`${API_BASE}/api/demo/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ wallet_address: randomAddress })
      });

      if (!response.ok) {
        throw new Error('failed to create demo user');
      }

      const payload = await response.json();
      const summary = payload.summary;

      setDemoMode(true);
      setDemoAddress(randomAddress);
      localStorage.setItem('demo_mode_active', 'true');
      localStorage.setItem('demo_wallet_address', randomAddress);

      if (summary) {
        applySummaryToState(summary);
      } else {
        localStorage.setItem('demo_data', JSON.stringify({ ...demoData }));
      }

      await refreshDemoFromServer(randomAddress);
    } catch (err) {
      console.error('Failed to enable demo mode', err);
      throw err;
    }
  };

  const disableDemoMode = async () => {
    if (demoAddress) {
      try {
        await fetch(`${API_BASE}/api/demo/exit`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ wallet_address: demoAddress })
        });
      } catch (err) {
        console.warn('Failed to notify backend about exiting demo mode', err);
      }
    }

    setDemoMode(false);
    setDemoAddress(null);
    setDemoData({ ...DEFAULT_DEMO_DATA });
    localStorage.removeItem('demo_mode_active');
    localStorage.removeItem('demo_wallet_address');
    localStorage.removeItem('demo_data');
  };

  const updateDemoData = (updates) => {
    setDemoData(prev => {
      const next = { ...prev, ...updates };
      localStorage.setItem('demo_data', JSON.stringify(next));
      return next;
    });
    setUpdateKey(prev => prev + 1); // 触发组件刷新
  };

  // 演示模式下的操作函数 - 积分兑换资产
  const demoDeposit = (protocol, amount) => {
    if (amount > demoData.points) {
      return { success: false, error: '积分不足' };
    }

    // 根据协议设置不同的APY和兑换率
    const protocolConfig = {
      'Uniswap V3': { apy: 45.2, exchangeRate: 0.01, asset: 'USDC' },
      'Aave V3': { apy: 8.5, exchangeRate: 0.01, asset: 'aUSDC' },
      'LoyaltyUSD': { apy: 0, exchangeRate: 0.01, asset: 'LUSD' },
      'LoyaltyX Protocol': { apy: 125, exchangeRate: 1, asset: 'Points' }
    };

    const config = protocolConfig[protocol] || { apy: 10, exchangeRate: 0.01, asset: 'USDC' };
    const exchangedAmount = amount * config.exchangeRate; // 兑换的资产数量

    const newDeposit = {
      protocol,
      amount, // 消耗的积分数量
      exchangedAmount, // 兑换得到的资产数量
      asset: config.asset,
      apy: config.apy,
      earned: 0,
      depositTime: Date.now()
    };

    updateDemoData({
      points: demoData.points - amount,
      defiDeposits: [...demoData.defiDeposits, newDeposit]
    });

    return { success: true, message: `成功用 ${amount} 积分兑换 ${exchangedAmount.toFixed(2)} ${config.asset}` };
  };

  const demoWithdraw = (index) => {
    const deposit = demoData.defiDeposits[index];
    if (!deposit) {
      return { success: false, error: '无效的存款' };
    }

    // 计算实时收益
    const actualEarned = calculateEarnings(deposit);
    const totalReturn = deposit.amount + actualEarned;
    const newDeposits = demoData.defiDeposits.filter((_, i) => i !== index);

    updateDemoData({
      points: demoData.points + totalReturn,
      defiDeposits: newDeposits
    });

    return { success: true, message: `成功取出 ${totalReturn.toFixed(2)} 积分（本金 + 收益）` };
  };

  // 稳定币铸造：100 积分 = 1 LUSD，需要 150% 抵押率
  const demoMintStablecoin = (lusdAmount) => {
    const mintRatio = 100; // 100 积分 = 1 LUSD
    const collateralRatio = 1.5; // 150% 抵押率

    // 计算需要的积分抵押：铸造 1 LUSD 需要 100 * 1.5 = 150 积分
    const requiredCollateral = lusdAmount * mintRatio * collateralRatio;

    if (requiredCollateral > demoData.points) {
      return { success: false, error: `需要 ${requiredCollateral.toFixed(0)} 积分作为抵押（当前仅有 ${demoData.points.toFixed(0)} 积分）` };
    }

    updateDemoData({
      points: demoData.points - requiredCollateral,
      collateral: demoData.collateral + requiredCollateral,
      stablecoinMinted: demoData.stablecoinMinted + lusdAmount
    });

    return { success: true, message: `成功铸造 ${lusdAmount} LUSD（消耗 ${requiredCollateral.toFixed(0)} 积分抵押）` };
  };

  const demoRedeemStablecoin = (amount) => {
    if (amount > demoData.stablecoinMinted) {
      return { success: false, error: 'LUSD 不足' };
    }

    const collateralReturn = (amount / demoData.stablecoinMinted) * demoData.collateral;

    updateDemoData({
      points: demoData.points + collateralReturn,
      collateral: demoData.collateral - collateralReturn,
      stablecoinMinted: demoData.stablecoinMinted - amount
    });

    return { success: true, message: `成功赎回 ${collateralReturn.toFixed(2)} 积分抵押物` };
  };

  const demoStake = (amount) => {
    if (amount > demoData.pfiTokens) {
      return { success: false, error: 'PFI 代币不足' };
    }

    updateDemoData({
      pfiTokens: demoData.pfiTokens - amount,
      stakedTokens: demoData.stakedTokens + amount
    });

    return { success: true, message: `成功质押 ${amount} PFI 代币` };
  };

  const demoUnstake = (amount) => {
    if (amount > demoData.stakedTokens) {
      return { success: false, error: '质押数量不足' };
    }

    // When unstaking, also claim accumulated staking rewards
    const rewards = demoData.stakingRewards;

    updateDemoData({
      pfiTokens: demoData.pfiTokens + amount,
      stakedTokens: demoData.stakedTokens - amount,
      points: demoData.points + rewards,
      stakingRewards: 0
    });

    return { success: true, message: `成功解除质押 ${amount} PFI 代币，并领取 ${rewards.toFixed(2)} 积分奖励` };
  };

  const value = {
    demoMode,
    demoAddress,
    demoData: {
      ...demoData,
      defiDeposits: getDepositsWithEarnings() // 返回带实时收益的数据
    },
    updateKey, // 暴露更新键
    enableDemoMode,
    disableDemoMode,
    updateDemoData,
    demoDeposit,
    demoWithdraw,
    demoMintStablecoin,
    demoRedeemStablecoin,
    demoStake,
    demoUnstake,
  };

  return (
    <DemoModeContext.Provider value={value}>
      {children}
    </DemoModeContext.Provider>
  );
}

export function useDemoMode() {
  const context = useContext(DemoModeContext);
  if (!context) {
    throw new Error('useDemoMode must be used within DemoModeProvider');
  }
  return context;
}
