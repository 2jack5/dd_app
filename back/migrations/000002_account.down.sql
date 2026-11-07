-- 回滚：用户、登录身份与协议
BEGIN;
DROP TABLE agreement_acceptances;
DROP TABLE auth_identities;
DROP TABLE users;
COMMIT;
