UPDATE transaction_outputs
SET
  spent = TRUE,
  spent_tx_id = ?
WHERE
  id IN (?)