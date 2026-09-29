create table user (
    id bigint unsigned not null auto_increment primary key comment "用户ID 主键",
    name varchar(50)                            not null comment "用户名",
    password varchar(100)                       not null comment "用户密码",
    phone varchar(100)                          not null default '' comment "用户手机号",
    avatar varchar(100)                         not null default '' comment "用户头像",
    email varchar(100)                          not null default '' comment "用户邮箱",
    role enum('admin','user','unknow')          not null default 'user' comment "用户身份",
    status enum('approved','blocked','unknow')  not null default 'unknow' comment "用户状态 approved已通过 blocked已拉黑",
    sex enum('male','female','secret','unknow') not null default 'unknow' comment "用户性别",
    birthday date                               default null comment "用户生日",
    signature varchar(255)                      default '' comment "用户主页简介",
    created_at datetime                         not null default current_timestamp comment "创建时间",
    updated_at datetime                         not null default current_timestamp on update current_timestamp comment "更新时间",
    deleted_at datetime                         default null comment "删除时间 软删除",

    unique key `uk_name` (name),
    key `idx_email` (email),
    key `idx_phone` (phone),
    key `idx_deleted` (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment "用户表"