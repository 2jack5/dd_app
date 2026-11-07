-- 基础扩展与共享审计函数
-- 按文件编号升序执行；本文件与对应 down 文件由迁移工具记录版本。
BEGIN;

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE FUNCTION dd_reject_event_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'table % is append-only', TG_TABLE_NAME USING ERRCODE = '55000';
END;
$$;

COMMIT;
