--
-- DB migration file
--

BEGIN;

ALTER TABLE whitelist
    RENAME COLUMN comment TO name;

COMMIT;

--
-- end of file
--