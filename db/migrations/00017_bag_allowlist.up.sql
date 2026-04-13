--
-- DB migration file
--

BEGIN;

-- create bag_allowlist table
CREATE TABLE bag_allowlist (
   id            serial PRIMARY KEY,
   name          TEXT NOT NULL DEFAULT '',
   comment       TEXT NOT NULL DEFAULT '',
   created_at    timestamp DEFAULT NOW()
);

-- create the unique bag_allowlist index
CREATE UNIQUE INDEX bag_allowlist_name_distinct_idx ON bag_allowlist(name);

-- add the bags we know about
INSERT INTO bag_allowlist(name) VALUES('virginia.edu/virginia.edu.libraETD');

COMMIT;

--
-- end of file
--