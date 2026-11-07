-- 回滚：陪练入驻与认证
BEGIN;
DROP TABLE coach_skills;
DROP TABLE identity_checks;
DROP TABLE coach_applications;
DROP TABLE coaches;
COMMIT;
