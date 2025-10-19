import React, { createContext, useContext, useEffect, useState } from "react";
import { getInjectedProvider } from "./provider";
import { BrowserProvider } from "ethers";

const WalletCtx = createContext(null);

export function WalletProvider({ children }) {
  const [address, setAddress] = useState(null);
  const [chainId, setChainId] = useState(null);
  const [signer, setSigner] = useState(null);

  const readAccounts = async () => {
    const p = getInjectedProvider();
    if (!p) return;
    const accs = await p.request({ method: "eth_accounts" });
    const cid = await p.request({ method: "eth_chainId" });
    setAddress(accs?.[0] ?? null);
    setChainId(cid ?? null);

    // Create ethers signer
    if (accs?.[0]) {
      try {
        const provider = new BrowserProvider(p);
        const ethersSigner = await provider.getSigner();
        setSigner(ethersSigner);
      } catch (e) {
        console.error("Failed to create signer:", e);
        setSigner(null);
      }
    } else {
      setSigner(null);
    }
  };

  useEffect(() => {
    readAccounts();
    const p = getInjectedProvider();
    if (!p) return;
    const onAcc = (accs) => setAddress(accs?.[0] ?? null);
    const onChn = (cid) => setChainId(cid ?? null);
    p.on?.('accountsChanged', onAcc);
    p.on?.('chainChanged', onChn);
    return () => {
      p.removeListener?.('accountsChanged', onAcc);
      p.removeListener?.('chainChanged', onChn);
    };
  }, []);

  const connect = async () => {
    const p = getInjectedProvider();
    if (!p) { alert('未检测到注入钱包，请安装/启用 MetaMask'); return; }
    await p.request({ method: 'eth_requestAccounts' });
    await readAccounts();
  };

  const disconnect = () => setAddress(null);

  const switchToSepolia = async () => {
    const p = getInjectedProvider(); if (!p) return;
    try {
      await p.request({ method: 'wallet_switchEthereumChain', params: [{ chainId: '0xaa36a7' }] });
    } catch (e) {
      if (e?.code === 4902) {
        await p.request({
          method: 'wallet_addEthereumChain',
          params: [{ chainId: '0xaa36a7', chainName: 'Sepolia', rpcUrls: ['https://rpc.sepolia.org'], nativeCurrency: { name: 'SepoliaETH', symbol: 'ETH', decimals: 18 } }]
        });
      } else { console.error(e); }
    }
  };

  return <WalletCtx.Provider value={{ address, chainId, signer, connect, disconnect, switchToSepolia }}>{children}</WalletCtx.Provider>;
}

export const useWallet = () => useContext(WalletCtx);
