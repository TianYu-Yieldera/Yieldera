# Airdrop Examples

This directory contains example files for the airdrop feature.

## Whitelist CSV Format

The whitelist CSV file should have the following format:

```csv
address,amount
0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1,1000
0x5B38Da6a701c568545dCfcB03FcB875f56beddC4,500
```

### Fields:
- `address`: Ethereum address (must start with 0x, 42 characters total)
- `amount`: Amount to allocate (decimal number, can include decimals like 100.5)

### Requirements:
- First row must be the header: `address,amount`
- Each address should be unique
- Addresses are case-insensitive (will be normalized to lowercase)
- Amount should be a valid number

## Admin Whitelist SQL

To add admin addresses to the whitelist, connect to PostgreSQL and run:

```sql
-- Add admin addresses
INSERT INTO admin_whitelist (address, name) VALUES
  ('0x742d35cc6634c0532925a3b844bc9e7595f0beb1', 'Admin 1'),
  ('0x5b38da6a701c568545dcfcb03fcb875f56beddc4', 'Admin 2')
ON CONFLICT (address) DO NOTHING;
```

## Testing Flow

### 1. Add yourself as admin
```bash
# Connect to PostgreSQL
docker exec -it loyalty-points-system-final-postgres-1 psql -U postgres -d loyalty_points

# Add your address
INSERT INTO admin_whitelist (address, name) VALUES
  ('YOUR_WALLET_ADDRESS_IN_LOWERCASE', 'Your Name');
```

### 2. Create a campaign via API

```bash
# Get auth token (use your wallet address as Bearer token for now)
curl -X POST http://localhost:8080/api/admin/airdrop/campaigns \
  -H "Authorization: Bearer YOUR_WALLET_ADDRESS" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Airdrop Campaign",
    "description": "Testing the airdrop feature",
    "start_time": "2025-01-01T00:00:00Z",
    "end_time": "2025-12-31T23:59:59Z",
    "total_budget": "10000",
    "is_demo": true
  }'
```

### 3. Import whitelist

```bash
# Upload CSV file (replace {CAMPAIGN_ID} with the ID from step 2)
curl -X POST http://localhost:8080/api/admin/airdrop/campaigns/{CAMPAIGN_ID}/allocations/import \
  -H "Authorization: Bearer YOUR_WALLET_ADDRESS" \
  -F "file=@examples/airdrop/whitelist-example.csv"
```

### 4. Activate campaign

```bash
curl -X POST http://localhost:8080/api/admin/airdrop/campaigns/{CAMPAIGN_ID}/activate \
  -H "Authorization: Bearer YOUR_WALLET_ADDRESS"
```

### 5. Check eligibility (anyone can do this)

```bash
curl "http://localhost:8080/api/airdrop/campaigns/{CAMPAIGN_ID}/eligibility?address=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"
```

### 6. Claim airdrop (requires wallet signature)

```javascript
// Frontend code example
const message = `Claim airdrop from campaign ${campaignId} with nonce ${nonce}`;
const signature = await signer.signMessage(message);

fetch(`http://localhost:8080/api/airdrop/campaigns/${campaignId}/claim`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    address: userAddress,
    nonce: nonce,  // Generate random nonce: Date.now().toString()
    signature: signature
  })
});
```

## Campaign Status Flow

1. **draft** - Created but not active yet (can edit/import allocations)
2. **scheduled** - Activated but start_time hasn't arrived
3. **active** - Currently active, users can claim
4. **claimable** - End time passed but still allowing claims (7 days)
5. **closed** - No longer accepting claims
6. **archived** - Historical record

The scheduler automatically transitions statuses based on time.
