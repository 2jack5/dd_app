-- 回滚：文件资源与通知
BEGIN;
DROP TABLE notifications;
DROP TABLE media_assets;
COMMIT;
