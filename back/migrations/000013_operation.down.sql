-- 回滚：首页运营、内容审核与处罚
BEGIN;
DROP TABLE user_penalties;
DROP TABLE content_audits;
DROP TABLE coach_recommendations;
DROP TABLE banners;
DROP TABLE activities;
COMMIT;
