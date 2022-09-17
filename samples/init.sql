USE cocoon;
DROP TABLE IF EXISTS samples;
CREATE TABLE samples(
    id INT PRIMARY KEY AUTO_INCREMENT,
    tiny_type TINYINT NOT NULL,
    int_type INT NOT NULL,
    big_int_type BIGINT NOT NULL,
    float_type FLOAT NOT NULL,
    double_type DOUBLE NOT NULL,
    decimal_type DECIMAL(8, 4) NOT NULL,
    date_type DATE NOT NULL,
    time_type TIME NOT NULL,
    year_type YEAR NOT NULL,
    datetime_type DATETIME NOT NULL,
    timestamp_type TIMESTAMP NOT NULL,
    char_type CHAR NOT NULL,
    varchar_type VARCHAR(128) NOT NULL,
    tinyblob_type TINYBLOB NOT NULL,
    tinytext_type TINYTEXT NOT NULL,
    blob_type BLOB NOT NULL,
    text_type TEXT NOT NULL,
    null_type INT
);