-- Create books table
CREATE TABLE books (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR (255)   NOT NULL default '',
    author      VARCHAR (255)   NOT NULL default '',
    content     text            NOT NULL default '',
    price       NUMERIC(15,2)   not null default 0,
    created_at  INT             not null default 0,
    updated_at  INT             not null default 0,
    rating      INT             not null default 0
);
CREATE INDEX idx_books_title ON books (title);
