-- 回滚：服务商品、排期与预约占位
BEGIN;
DROP TABLE booking_holds;
DROP TABLE coach_availability;
DROP TABLE coach_service_regions;
DROP TABLE coach_services;
COMMIT;
