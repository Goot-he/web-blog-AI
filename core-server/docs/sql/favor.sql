create table interaction_favor(
    id bigint unsigned                  not null auto_increment primary key comment "收藏记录ID",
    created_at datetime                 not null default current_timestamp comment "创建时间",
    updated_at datetime                 not null default current_timestamp on update current_timestamp comment "更新时间",
    user_id bigint unsigned             not null default 0 comment "收藏用户ID",
    object_type enum('article','unknow') not null default 'unknow' comment "收藏对象 文章",
    object_id bigint unsigned           not null default 0 comment "收藏对象ID",
    status enum('collected','nothing')   not null default 'nothing' comment "收藏状态 collected已收藏 nothing已取消",

    unique key `uk_user_object_type_id` (user_id,object_type,object_id),
    key `idx_user_id` (user_id),
    key `idx_created_at` (created_at)
) ENGINE = InnoDB charset = utf8mb4 collate = utf8mb4_general_ci comment "收藏表"