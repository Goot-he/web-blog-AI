create table interaction_like(
    id bigint unsigned                  not null auto_increment primary key comment "点赞记录ID",
    created_at datetime                 not null default current_timestamp comment "创建时间",
    updated_at datetime                 not null default current_timestamp on update current_timestamp comment "更新时间",
    user_id bigint unsigned             not null default 0 comment "点赞用户ID",
    object_type enum('article','comment','unknow') not null default 'unknow' comment "点赞对象 评论、文章",
    object_id bigint unsigned           not null default 0 comment "点赞对象ID",
    status enum('thumb_up','nothing')   not null default 'nothing' comment "点赞状态 thumb_up已点赞 nothing已取消",

    unique key `uk_user_object_type_id` (user_id,object_type,object_id),
    key `idx_user_id` (user_id),
    key `idx_created_at` (created_at)
) ENGINE = InnoDB charset = utf8mb4 collate = utf8mb4_general_ci comment "点赞表"