--
-- DB migration file
--

BEGIN;

-- create submission state table
CREATE TABLE submission_state (
   id            serial PRIMARY KEY,
   submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
   status        submission_states NOT NULL DEFAULT 'registered',
   created_at    timestamp DEFAULT NOW()
);

-- create the submission_state index(s)
CREATE INDEX submission_state_submission_key_idx ON submission_state(submission);

-- create bag state table
CREATE TABLE bag_state (
   id            serial PRIMARY KEY,
   submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
   name          VARCHAR( 64 ) NOT NULL DEFAULT '',
   status        bag_states NOT NULL DEFAULT 'registered',
   created_at    timestamp DEFAULT NOW()
);

-- create the bag_state index(s)
CREATE INDEX bag_state_submission_name_key_idx ON bag_state(submission, name);

-- drop the status fields from the content tables
ALTER TABLE submissions
   DROP COLUMN status,
   DROP COLUMN updated_at;

ALTER TABLE bags
   DROP COLUMN status,
   DROP COLUMN updated_at;

-- auto vacuum parameters
-- see: https://aws.amazon.com/blogs/database/understanding-autovacuum-in-amazon-rds-for-postgresql-environments/

ALTER TABLE submission_state SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE submission_state SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE submission_state SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE submission_state SET (autovacuum_analyze_threshold = 1000);

ALTER TABLE bag_state SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE bag_state SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE bag_state SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE bag_state SET (autovacuum_analyze_threshold = 1000);

COMMIT;

--
-- end of file
--