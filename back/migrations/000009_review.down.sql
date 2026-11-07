-- 回滚：评价与内容版本
BEGIN;
DROP TABLE review_contents;
DROP TABLE reviews;
COMMIT;
