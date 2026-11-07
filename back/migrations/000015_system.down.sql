-- 回滚：业务规则、幂等与可靠任务
BEGIN;
DROP TABLE outbox_events;
DROP TABLE idempotency_keys;
DROP TABLE business_rule_versions;
COMMIT;
