--
-- DB migration file
--

BEGIN;

ALTER TABLE clients
    DROP COLUMN approval_email;

COMMIT;

--
-- end of file
--