--
-- DB migration file
--

BEGIN;

ALTER TABLE submissions
    ADD COLUMN collection_name VARCHAR( 64 ) NOT NULL DEFAULT '';

COMMIT;

--
-- end of file
--