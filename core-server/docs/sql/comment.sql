create table comment (
    id bigint unsigned not null auto_increment primary key comment "评论ID",
    created_at datetime         not null default current_timestamp comment "创建时间",
    deleted_at datetime         default null comment "删除时间 软删除",
    article_id bigint unsigned  not null default 0 comment "评论文章ID",
    user_id bigint unsigned     not null default 0 comment "评论用户ID",
    parent_id bigint unsigned   not null default 0 comment "父评论ID",
    root_id bigint unsigned     not null default 0 comment "根评论ID",
    reply_to_id bigint unsigned not null default 0 comment "回复用户ID",
    content text                not null comment "评论内容",
    like_count bigint unsigned  not null default 0 comment "评论点赞量",
    child_count bigint unsigned not null default 0 comment "子评论数量",
    is_top tinyint(1)           not null default false comment "是否置顶评论",

    index `idx_article_id` (article_id),
    index `idx_user_id` (user_id),
    index `idx_root_id` (root_id),
    index `idx_deleted_at` (deleted_at)
) ENGINE = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci comment "评论表"