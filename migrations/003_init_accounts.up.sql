CREATE TABLE accounts (
    id          INT UNSIGNED    NOT NULL AUTO_INCREMENT COMMENT '流水號',
    userID      INT UNSIGNED    NOT NULL COMMENT '使用者ID',
    bankID      INT UNSIGNED    NOT NULL COMMENT '銀行ID',
    accountName VARCHAR(100)    NOT NULL COMMENT '帳戶名稱',
    balance     BIGINT          NOT NULL DEFAULT 0 COMMENT '帳戶餘額',
    createdAt   DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt   DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX (userID),
    INDEX (bankID),
    FOREIGN KEY (userID) REFERENCES users(id),
    FOREIGN KEY (bankID) REFERENCES banks(id)
);