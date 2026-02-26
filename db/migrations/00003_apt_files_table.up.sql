--
-- DB migration file
--

BEGIN;

CREATE TABLE apt_files (
   id            serial PRIMARY KEY,
   file_name     TEXT NOT NULL DEFAULT '',
   hash          VARCHAR( 64 ) NOT NULL DEFAULT '',
   bag_name      VARCHAR( 64 ) NOT NULL DEFAULT '',
   apt_added_at  timestamp DEFAULT NOW(),
   created_at    timestamp DEFAULT NOW()
);

-- create the apt_files index(s)
CREATE INDEX apt_files_hash_key_idx ON apt_files(hash);
CREATE INDEX apt_files_hash_bag_key_idx ON apt_files(bag_name);

ALTER TABLE apt_files SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE apt_files SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE apt_files SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE apt_files SET (autovacuum_analyze_threshold = 1000);

COMMIT;

--
-- end of file
--