-- 回滚：售后工单与工单事件
BEGIN;
DROP TABLE case_events;
DROP TABLE after_sales_cases;
COMMIT;
