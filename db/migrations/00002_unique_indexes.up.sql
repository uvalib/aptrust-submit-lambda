--
-- DB migration file
--

BEGIN;

-- create the unique bag index
CREATE UNIQUE INDEX bags_submission_name_distinct_idx ON bags(submission, name);

-- create the unique file index
CREATE UNIQUE INDEX files_submission_bag_name_name_distinct_idx ON files(submission, bag_name, name);

COMMIT;

--
-- end of file
--