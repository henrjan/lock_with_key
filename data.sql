-- SQLite
SELECT hst1.owner_id, acc.account_name, acc.account_number, hst1.current_balance 
FROM history_logs hst1
INNER JOIN (
    SELECT owner_id, MAX(id) as id 
    FROM history_logs hst
    GROUP BY hst.owner_id
) hst2 ON hst1.owner_id = hst2.owner_id and hst1.id = hst2.id
INNER JOIN accounts acc ON hst1.owner_id = acc.id
ORDER BY hst1.owner_id;

SELECT hst.id, hst.owner_id, acc.account_name, acc.account_number, hst.last_balance, hst.balance, hst.current_balance 
FROM history_logs hst
INNER JOIN accounts acc ON hst.owner_id = acc.id
WHERE hst.last_balance <> 0 and hst.owner_id = 10
ORDER BY hst.id;