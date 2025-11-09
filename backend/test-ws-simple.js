const WebSocket = require('ws');

const ws = new WebSocket('wss://arb-sepolia.g.alchemy.com/v2/FP6JOVxZoc4lDScODskcP');

ws.on('open', () => {
  console.log('✅ WebSocket connected!');
  ws.send(JSON.stringify({
    jsonrpc: '2.0',
    id: 1,
    method: 'eth_blockNumber',
    params: []
  }));
});

ws.on('message', (data) => {
  console.log('✅ Received:', data.toString());
  ws.close();
});

ws.on('error', (error) => {
  console.error('❌ WebSocket error:', error.message);
  console.error('Full error:', error);
});

ws.on('close', (code, reason) => {
  console.log('WebSocket closed:', code, reason.toString());
  process.exit(0);
});

setTimeout(() => {
  console.log('⏱️ Timeout - no response');
  process.exit(1);
}, 10000);
