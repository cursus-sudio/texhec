CREATE TABLE saves(
    id UUID PRIMARY KEY,
    created DATETIME NOT NULL,
    last_modified DATETIME NOT NULL,
    name VARCHAR(64) NOT NULL
);
