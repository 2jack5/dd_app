-- 回滚：后台权限、复核与审计
BEGIN;
DROP TABLE admin_audit_logs;
DROP TABLE admin_approvals;
DROP TABLE admin_role_permissions;
DROP TABLE admin_user_roles;
DROP TABLE admin_permissions;
DROP TABLE admin_roles;
DROP TABLE admin_users;
COMMIT;
