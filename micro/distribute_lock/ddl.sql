create table user_balance
(
    `uid`     int(11) NOT NULL AUTO_INCREMENT COMMENT '用户uid',
    `name`    varchar(255) NOT NULL DEFAUlT '' COMMENT '用户名',
    `balance` int(11) NOT NULL DEFAULT 0 COMMENT '用户余额（现金）',
    `version` int(11) NOT NULL DEFAULT 0 COMMENT '余额修改版本号',
    PRIMARY KEY (`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户余额表';