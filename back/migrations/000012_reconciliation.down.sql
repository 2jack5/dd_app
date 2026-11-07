-- 回滚：渠道对账与差异
BEGIN;
DROP TABLE reconciliation_diffs;
DROP TABLE reconciliation_lines;
DROP TABLE reconciliation_runs;
COMMIT;
