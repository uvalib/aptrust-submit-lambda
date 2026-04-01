--
-- DB migration file
--

BEGIN;

-- create submission_failures table
CREATE TABLE submission_failures (
   id               serial PRIMARY KEY,
   submission       VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
   created_at       timestamp DEFAULT NOW()
);

-- create the submission_failures index(s)
CREATE INDEX submission_failures_submission_idx ON submission_failures(submission);

-- auto vacuum parameters
-- see: https://aws.amazon.com/blogs/database/understanding-autovacuum-in-amazon-rds-for-postgresql-environments/

ALTER TABLE submission_failures SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE submission_failures SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE submission_failures SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE submission_failures SET (autovacuum_analyze_threshold = 1000);

COMMIT;

--
-- end of file
--