-- Set timezone
-- Example: Asia/Jakarta
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="UTC";

-- UUID extension for universally unique identifiers (UUIDs)
-- https://www.postgresql.org/docs/current/uuid-ossp.html
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";