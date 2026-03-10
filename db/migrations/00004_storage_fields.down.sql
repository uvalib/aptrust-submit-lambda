--
-- DB migration file
--

BEGIN;

ALTER TABLE clients
    DROP COLUMN default_storage;

ALTER TABLE submissions
    DROP COLUMN storage;

COMMIT;

--
-- end of file
--