create table category (
    id bigint unsigned          not null auto_increment primary key comment "分类ID",
    created_at datetime         not null default current_timestamp comment "创建时间",
    updated_at datetime         not null default current_timestamp on update current_timestamp comment "更新时间",
    deleted_at datetime         default null comment "删除时间 软删除",

    parent_id bigint unsigned   not null default 0 comment "一级分类ID",
    name varchar(50)            not null comment "分类名称",

    unique key `uk_parent_id_name` (parent_id,name),
    key `idx_category_parent_id` (parent_id)
) ENGINE=InnoDB default charset = utf8mb4 comment "分类表";

