CREATE TABLE banks (
    id              INT UNSIGNED        NOT NULL AUTO_INCREMENT COMMENT '流水號',
    code            VARCHAR(10)         NOT NULL COMMENT '銀行代號',
    bankName        VARCHAR(100)        NOT NULL COMMENT '銀行名稱',
    capitalAmount   INT UNSIGNED        NOT NULL DEFAULT 0 COMMENT '資本額',
    createdAt       DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt       DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (code)
);
