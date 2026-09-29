create table article (
    id bigint unsigned auto_increment primary key comment "文章ID",

    created_at      datetime        not null default current_timestamp comment "创建时间",
    updated_at      datetime        not null default current_timestamp on update current_timestamp comment "更新时间",
    deleted_at      datetime        not null comment "删除时间 软删除",
    published_at    datetime        not null comment "发布时间 可为空 表示未发布",

    title           varchar(100)        not null comment "文章标题",
    summary         varchar(500)        not null default '' comment "文章简述",
    content         longtext            not null comment "文章正文",
    cover_image     varchar(500)        not null default '' comment "文章封面URL",
    auth_id         bigint unsigned     not null comment "作者id",
    category_id     bigint unsigned     not null default 0 comment "分类id",

    is_top          tinyint(1)          not null default false comment "是否置顶 默认为否",
    is_published    tinyint(1)          not null default true comment "是否发布 默认发布",

    favor_count     bigint unsigned     not null default 0 comment "收藏数量",
    comment_count   bigint unsigned     not null default 0 comment "评论数量",
    like_count      bigint unsigned     not null default 0 comment "喜欢数量",
    review_count    bigint unsigned     not null default 0 comment "浏览数量"
) comment "博客文章表";

create index idx_category_id on article(category_id);
create index idx_auth_id on article(auth_id)