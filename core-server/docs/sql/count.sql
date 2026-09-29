create table interaction_count (
    id bigint unique            not null auto_increment primary key comment "交互记录ID",
    created_at datetime         not null default current_timestamp comment "创建时间",
    updated_at datetime         not null default current_timestamp on update current_timestamp comment "更新时间",
    object_type varchar(50)     not null comment "对象类型 譬如文章",
    object_id bigint unsigned   not null comment "对象ID 譬如文章ID",
    interaction_type enum('like','favor','comment','unknow')
                                not null default 'unknow' comment "交互类型 点赞、收藏、评论",
    count bigint unsigned       not null default 0 comment "计数",
    constraint uk_object unique (object_type,object_id,interaction_type)
) comment "交互数据统计表";
