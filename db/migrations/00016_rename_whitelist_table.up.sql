--
-- DB migration file
--

BEGIN;

-- rename the table
ALTER TABLE whitelist RENAME TO hash_allowlist;

-- drop the previous index(es)
DROP INDEX whitelist_hash_distinct_idx;

-- create the hash_allowlist index(es)
CREATE UNIQUE INDEX hash_allowlist_hash_distinct_idx ON hash_allowlist(hash);

COMMIT;

--
-- end of file
--