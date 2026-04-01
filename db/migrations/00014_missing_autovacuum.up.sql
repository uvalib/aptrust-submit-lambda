--
-- DB migration file
--

BEGIN;

-- auto vacuum parameters
-- see: https://aws.amazon.com/blogs/database/understanding-autovacuum-in-amazon-rds-for-postgresql-environments/

ALTER TABLE approvals SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE approvals SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE approvals SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE approvals SET (autovacuum_analyze_threshold = 1000);

ALTER TABLE submission_conflicts SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE submission_conflicts SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE submission_conflicts SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE submission_conflicts SET (autovacuum_analyze_threshold = 1000);

COMMIT;

--
-- end of file
--