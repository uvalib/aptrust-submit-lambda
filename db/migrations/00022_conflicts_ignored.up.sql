--
-- DB migration file
--

BEGIN;

ALTER TABLE submission_conflicts
    ADD COLUMN ignored BOOLEAN DEFAULT FALSE;

COMMIT;

--
-- end of file
--