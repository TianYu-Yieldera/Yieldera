export function getInjectedProvider() {
  if (typeof window === "undefined") return null;
  const eth = window.ethereum;
  if (eth?.providers?.length) {
    const mm = eth.providers.find((p) => p.isMetaMask);
    return mm || eth.providers[0];
  }
  if (eth) return eth;
  let selected = null;
  function onAnnounce(e){ const p = e?.detail?.provider; if (p?.isMetaMask) selected=p; if(!selected && p) selected=p; }
  window.addEventListener('eip6963:announceProvider', onAnnounce, { once:true });
  window.dispatchEvent(new Event('eip6963:requestProvider'));
  return selected;
}
