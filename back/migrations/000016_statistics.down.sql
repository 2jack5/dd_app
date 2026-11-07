-- 回滚：活跃事件与统计
BEGIN;
DROP TABLE daily_metrics;
DROP TABLE user_activity_events;
COMMIT;
