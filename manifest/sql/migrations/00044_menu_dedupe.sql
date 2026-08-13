-- 菜单去重(修复侧栏出现两个「视频管理」)。
--
-- 成因: 00042/00043 清理旧顶级目录时是按写死的 id 删的(id=67 等), 但如果库里的
-- 菜单行不是全部来自迁移(比如早期在「菜单权限」页面手工建过节点), id 就对不上,
-- 旧节点删不掉, 于是新旧两个「视频管理」同时挂在顶层。
--
-- 修法改为**按 (父级, 名称) 去重**, 不再依赖 id:
--   * 同一父级下同名的菜单节点只保留一个 —— 优先留 sort 小的(新结构的一级 sort 是
--     20~23, 旧顶级是 55 这类大值), sort 相同留 id 小的(迁移种子节点比手工建的早);
--   * 被删节点的子节点(二级页面、接口权限)先整体挂到保留节点下, 再删, 不丢权限数据;
--   * 合并后再对接口权限按 (父级, 路径, 方法) 去重。
-- 跑多轮直到无重复(合并一级后二级可能出现新的同名对)。本迁移幂等, 重复执行无害。
-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    keeper bigint;
    r record;
    d record;
BEGIN
    -- 最多三轮: 一级合并可能让二级出现重复, 三层树足够收敛
    FOR pass IN 1..3 LOOP
        FOR r IN
            SELECT parent_id, name FROM admin_permission
            WHERE is_menu = 1
            GROUP BY parent_id, name
            HAVING count(*) > 1
        LOOP
            SELECT id INTO keeper FROM admin_permission
            WHERE parent_id = r.parent_id AND name = r.name AND is_menu = 1
            ORDER BY sort ASC, id ASC
            LIMIT 1;
            FOR d IN
                SELECT id FROM admin_permission
                WHERE parent_id = r.parent_id AND name = r.name
                  AND is_menu = 1 AND id <> keeper
            LOOP
                UPDATE admin_permission SET parent_id = keeper WHERE parent_id = d.id;
                DELETE FROM admin_permission WHERE id = d.id;
            END LOOP;
        END LOOP;
    END LOOP;

    -- 接口权限(is_menu=0)按 (父级, 路径, 方法) 去重, 留最早的一条
    FOR r IN
        SELECT parent_id, route_url, method FROM admin_permission
        WHERE is_menu = 0
        GROUP BY parent_id, route_url, method
        HAVING count(*) > 1
    LOOP
        DELETE FROM admin_permission
        WHERE is_menu = 0 AND parent_id = r.parent_id
          AND route_url = r.route_url AND method = r.method
          AND id <> (SELECT min(id) FROM admin_permission
                     WHERE is_menu = 0 AND parent_id = r.parent_id
                       AND route_url = r.route_url AND method = r.method);
    END LOOP;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- 去重不可逆(被合并的重复节点无法还原), 也无需还原。
SELECT 1;
