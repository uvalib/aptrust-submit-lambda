--
-- DB migration file
--

BEGIN;

-- create files table
CREATE TABLE whitelist (
   id            serial PRIMARY KEY,
   name          TEXT NOT NULL DEFAULT '',
   hash          VARCHAR( 64 ) NOT NULL DEFAULT '',
   created_at    timestamp DEFAULT NOW()
);

-- create the unique bag index
CREATE UNIQUE INDEX whitelist_hash_distinct_idx ON whitelist(hash);

COMMIT;

--
-- end of file
--