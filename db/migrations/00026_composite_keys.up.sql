--
-- DB migration file
--

BEGIN;

ALTER TABLE bag_states
    DROP CONSTRAINT bag_state_submission_fkey;

ALTER TABLE files
    DROP CONSTRAINT files_submission_fkey;

ALTER TABLE bag_states
    ADD CONSTRAINT bag_states_submission_bag_name_fkey
    FOREIGN KEY (submission, bag_name)
    REFERENCES bags(submission, bag_name);

ALTER TABLE files
    ADD CONSTRAINT files_submission_bag_name_fkey
    FOREIGN KEY (submission, bag_name)
    REFERENCES bags(submission, bag_name);

COMMIT;

--
-- end of file
--