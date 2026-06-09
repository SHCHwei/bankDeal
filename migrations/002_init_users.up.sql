CREATE TABLE users (
    id          INT UNSIGNED        NOT NULL AUTO_INCREMENT COMMENT '流水號',
    firstName   VARCHAR(50)         NOT NULL COMMENT '名字',
    lastName    VARCHAR(50)         NOT NULL COMMENT '姓氏',
    email       VARCHAR(255)        NOT NULL COMMENT '電子郵件',
    phone       VARCHAR(20)         NOT NULL COMMENT '電話',
    birthDate   DATE                NULL COMMENT '生日',
    createdAt   DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt   DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (email)
);