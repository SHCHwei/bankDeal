CREATE TABLE deals (
    id                  INT UNSIGNED        NOT NULL AUTO_INCREMENT COMMENT '流水號',
    accountID           INT UNSIGNED        NOT NULL COMMENT '帳戶ID',
    volume              BIGINT              NOT NULL DEFAULT 0 COMMENT '交易金額',
    transactionType     TINYINT UNSIGNED    NOT NULL COMMENT '交易類型，0: 存款, 1: 提款',
    tradingAccountID    INT UNSIGNED        NOT NULL COMMENT '目標帳戶ID',
    remark              TEXT                COMMENT '備註',
    createdAt           DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tradingAccountID) REFERENCES accounts(id),
    FOREIGN KEY (accountID) REFERENCES accounts(id),
    INDEX (accountID),
    INDEX (tradingAccountID)
);
