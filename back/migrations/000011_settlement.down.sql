-- 回滚：结算账户、佣金、结算、追偿与补款
BEGIN;
DROP TABLE settlement_supplements;
DROP TABLE settlement_batch_items;
DROP TABLE settlement_batches;
DROP TABLE settlement_adjustments;
DROP TABLE settlement_items;
DROP TABLE commission_rules;
DROP TABLE settlement_accounts;
DROP FUNCTION dd_guard_settlement_supplement();
DROP FUNCTION dd_guard_settlement_adjustment();
COMMIT;
