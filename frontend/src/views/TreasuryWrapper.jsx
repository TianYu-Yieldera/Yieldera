/**
 * Treasury Wrapper
 *
 * 为所有 Treasury 相关页面提供 Base Smart Wallet 环境
 * Treasury 功能仅在 Base 链可用
 */

import React from 'react';
import { Outlet } from 'react-router-dom';
import BaseSmartWalletProvider from '../web3/BaseSmartWalletProvider';

export function TreasuryWrapper() {
  return (
    <BaseSmartWalletProvider>
      <Outlet />
    </BaseSmartWalletProvider>
  );
}

export default TreasuryWrapper;
