-- 回滚：基础扩展与共享审计函数
BEGIN;
DROP FUNCTION dd_reject_event_mutation();
-- btree_gist 可能由其他 schema 共用，回滚保留该扩展。
COMMIT;
