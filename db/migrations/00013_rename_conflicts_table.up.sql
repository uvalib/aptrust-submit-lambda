--
-- DB migration file
--

BEGIN;

-- rename the table
ALTER TABLE conflicts RENAME TO submission_conflicts;

-- drop the previous index(es)
DROP INDEX conflicts_submission_idx;

-- create the submission_conflicts index(es)
CREATE INDEX submission_conflicts_submission_idx ON submission_conflicts(submission);

COMMIT;

--
-- end of file
--