--
-- DB migration file
--

BEGIN;

-- create storage_options table
CREATE TABLE storage_options (
   id          serial PRIMARY KEY,
   value       VARCHAR( 64 ) NOT NULL DEFAULT '',
   label       VARCHAR( 64 ) NOT NULL DEFAULT '',
   is_active   BOOLEAN DEFAULT TRUE,
   created_at  timestamp DEFAULT NOW()
);

INSERT INTO storage_options(value, label) VALUES('Standard', 'Standard');
INSERT INTO storage_options(value, label) VALUES('Glacier-Deep-OH', 'Glacier-Deep-OH');
INSERT INTO storage_options(value, label) VALUES('Glacier-VA', 'Glacier-VA');

-- auto vacuum parameters
-- see: https://aws.amazon.com/blogs/database/understanding-autovacuum-in-amazon-rds-for-postgresql-environments/

ALTER TABLE storage_options SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE storage_options SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE storage_options SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE storage_options SET (autovacuum_analyze_threshold = 1000);

COMMIT;

--
-- end of file
--