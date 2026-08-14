-- 用户组对齐公司后台「会员等级」(Mongo user_group):
-- 售卖字段(价格/天数/封面/赠币等)原先拆在 vip_package, 后台会员等级页需要合在同一行展示。
-- 原 status(1启用/0停用) 对应公司 is_disabled(0正常/1禁用), 含义相反, 不改列以免破坏已有数据。
-- 原 remark 继续当「描述」用。
-- +goose Up
ALTER TABLE user_group
    ADD COLUMN IF NOT EXISTS img               varchar(512)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS title_heat        varchar(16)    NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS title_description varchar(64)    NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS title_picture     varchar(512)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS level             int            NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS promotion_type    int            NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS price             numeric(14,2)  NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS old_price         numeric(14,2)  NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS day_num           int            NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS gift_num          int            NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS download_num      int            NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS day_tips          varchar(128)   NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS price_tips        varchar(128)   NOT NULL DEFAULT '';

COMMENT ON COLUMN user_group.img IS '封面 URL';
COMMENT ON COLUMN user_group.title_heat IS '会员头衔, 空则前台默认学徒';
COMMENT ON COLUMN user_group.level IS '1普通 2普通+暗网';
COMMENT ON COLUMN user_group.promotion_type IS '0正常价格 1新人专享';
COMMENT ON COLUMN user_group.day_num IS '可用天数';
COMMENT ON COLUMN user_group.gift_num IS '赠送金币';
COMMENT ON COLUMN user_group.rate IS '购片折扣(%); -1金币视频免费 -2视频和帖子免费';

-- +goose Down
ALTER TABLE user_group
    DROP COLUMN IF EXISTS img,
    DROP COLUMN IF EXISTS title_heat,
    DROP COLUMN IF EXISTS title_description,
    DROP COLUMN IF EXISTS title_picture,
    DROP COLUMN IF EXISTS level,
    DROP COLUMN IF EXISTS promotion_type,
    DROP COLUMN IF EXISTS price,
    DROP COLUMN IF EXISTS old_price,
    DROP COLUMN IF EXISTS day_num,
    DROP COLUMN IF EXISTS gift_num,
    DROP COLUMN IF EXISTS download_num,
    DROP COLUMN IF EXISTS day_tips,
    DROP COLUMN IF EXISTS price_tips;
