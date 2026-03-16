--
-- DB migration file
--

BEGIN;

-- drop the table if it exists
DROP TABLE IF EXISTS submission_status;
DROP TABLE IF EXISTS bag_status;

ALTER TABLE submissions
    ADD COLUMN status submission_states NOT NULL DEFAULT 'registered',
    ADD COLUMN updated_at timestamp DEFAULT NOW();

ALTER TABLE bags
    ADD COLUMN status bag_states NOT NULL DEFAULT 'registered',
    ADD COLUMN updated_at timestamp DEFAULT NOW();

COMMIT;

--
-- end of file
--