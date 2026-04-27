--
-- DB migration file
--

BEGIN;

ALTER TABLE submission_conflicts
    DROP COLUMN ignored;

COMMIT;

--
-- end of file
--