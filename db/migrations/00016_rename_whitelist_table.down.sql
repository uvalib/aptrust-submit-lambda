--
-- DB migration file
--

BEGIN;

ALTER TABLE hash_allowlist RENAME TO whitelist;

DROP INDEX hash_allowlist_hash_distinct_idx;
CREATE UNIQUE INDEX whitelist_hash_distinct_idx ON whitelist(hash);

COMMIT;

--
-- end of file
--