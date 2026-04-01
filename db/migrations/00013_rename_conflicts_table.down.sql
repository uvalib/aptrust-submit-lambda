--
-- DB migration file
--

BEGIN;

ALTER TABLE submission_conflicts RENAME TO conflicts;

DROP INDEX submission_conflicts_submission_idx;
CREATE INDEX conflicts_submission_idx ON conflicts(submission);

COMMIT;

--
-- end of file
--