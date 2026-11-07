-- 回滚：支付、退款及资金事件
BEGIN;
DROP TABLE refund_events;
DROP TABLE refunds;
DROP TABLE payment_events;
DROP TABLE payments;
DROP FUNCTION dd_guard_payment_event();
COMMIT;
