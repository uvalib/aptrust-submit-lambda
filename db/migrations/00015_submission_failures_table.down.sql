--
-- DB migration file
--

BEGIN;

-- drop the index(s)
DROP INDEX IF EXISTS submission_failures_submission_idx;

-- drop the table if it exists
DROP TABLE IF EXISTS submission_failures;

COMMIT;

--
-- end of file
--