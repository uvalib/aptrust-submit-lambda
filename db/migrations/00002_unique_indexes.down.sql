--
-- DB migration file
--

BEGIN;

DROP INDEX IF EXISTS files_submission_bag_name_name_distinct_idx;
DROP INDEX IF EXISTS bags_submission_name_distinct_idx;

COMMIT;

--
-- end of file
--