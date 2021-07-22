-- Create email_validation table
CREATE TABLE email_validation (
    email       VARCHAR (255) NOT NULL default '',
    secret_code CHAR (6)      NOT NULL default '',

    PRIMARY KEY (email)
);

-- Comments
comment on column "email_validation".secret_code is '6 digits number';

-- Indexes
CREATE INDEX idx_email_validation_email_secret_code ON email_validation (email,secret_code);