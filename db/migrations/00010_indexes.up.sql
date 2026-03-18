--
-- DB migration file
--

BEGIN;

CREATE INDEX files_hash_key_idx ON files(hash);
CREATE INDEX apt_files_bag_name_key_idx ON apt_files(bag_name);
-- incorrect name
DROP INDEX apt_files_hash_bag_name_key_idx;

COMMIT;

--
-- end of file
--