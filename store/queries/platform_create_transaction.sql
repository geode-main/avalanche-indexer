INSERT INTO transactions (
  id,
  reference_tx_id,
  status,
  type,
  block,
  block_height,
  chain,
  memo,
  memo_text,
  fee,
  nonce,
  source_chain,
  destination_chain,
  timestamp,
  metadata
)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
ON CONFLICT (id) DO NOTHING