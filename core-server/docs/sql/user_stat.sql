create table user_status (
    id bigint unsigned                  not null auto_increment primary key comment "统计记录ID",
    user_id bigint unsigned             not null comment "用户ID",
    article_count bigint unsigned       not null default 0 comment "文章总数",
    followers_count bigint unsigned     not null default 0 comment "粉丝数",
    following_count bigint unsigned     not null default 0 comment "关注数",
    like_count bigint unsigned          not null default 0 comment "点赞数",
    receive_like_count bigint unsigned  not null default 0 comment "收到的点赞数",
    favor_count bigint unsigned         not null default 0 comment "收藏数",
    receive_favor_count bigint unsigned not null default 0 comment "收到的收藏数",
    created_at datetime                 not null default current_timestamp comment "创建时间",
    updated_at datetime                 not null default current_timestamp on update current_timestamp comment "更新时间",
    deleted_at datetime                 default null comment "删除时间 软删除",

    unique key `uk_user_id` (user_id),
    constraint `fk_user_stat_user_id` foreign key (`user_id`) references `user` (`id`)
) engine=InnoDB default charset = utf8mb4 collate =utf8mb4_general_ci comment "用户状态统计表"