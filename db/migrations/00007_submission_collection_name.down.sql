--
-- DB migration file
--

BEGIN;

ALTER TABLE submissions
    DROP COLUMN collection_name;

COMMIT;

--
-- end of file
--