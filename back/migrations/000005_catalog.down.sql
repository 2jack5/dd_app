-- 回滚：游戏、区服、段位与服务分类
BEGIN;
DROP TABLE service_categories;
DROP TABLE game_ranks;
DROP TABLE game_regions;
DROP TABLE games;
COMMIT;
