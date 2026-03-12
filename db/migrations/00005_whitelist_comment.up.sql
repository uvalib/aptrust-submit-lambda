--
-- DB migration file
--

BEGIN;

ALTER TABLE whitelist
    RENAME COLUMN name TO comment;

COMMIT;

--
-- end of file
--