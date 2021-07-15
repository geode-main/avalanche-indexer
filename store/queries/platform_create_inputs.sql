INSERT INTO transaction_inputs (id, tx_id)
VALUES @values
ON CONFLICT (id) DO NOTHING