--
-- DB migration file
--

BEGIN;

DROP INDEX files_hash_key_idx;
DROP INDEX apt_files_bag_name_key_idx;
CREATE INDEX apt_files_hash_bag_name_key_idx ON apt_files(bag_name);

COMMIT;

--
-- end of file
--