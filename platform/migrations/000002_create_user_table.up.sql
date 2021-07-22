-- create user table
create table "user"
(
    id                   uuid         not null default uuid_generate_v4() constraint user_pkey primary key,
    username             varchar(255) not null default '',
    full_name            varchar(255) not null default '',
    auth_key             varchar(32)  not null default '',
    password_hash        varchar(255) not null default '',
    password_reset_token varchar(255) not null constraint user_password_reset_token_key unique,
    verification_token   varchar(255) not null default '' constraint user_verification_token_key unique,
    email                varchar(255) not null default '',
    status               smallint     not null default 0,
    created_at           integer      not null default 0,
    updated_at           integer      not null default 0
);

comment on column "user".status is '0=Inactive, 1=Active, 2=Deleted';

CREATE UNIQUE INDEX user_index_email on "user" (LOWER(email));
CREATE UNIQUE INDEX user_index_username on "user" (LOWER(username));
CREATE INDEX idx_user_verification_token ON "user" (verification_token);
CREATE INDEX idx_user_password_reset_token ON "user"(password_reset_token);
CREATE INDEX idx_user_status ON "user"(status);