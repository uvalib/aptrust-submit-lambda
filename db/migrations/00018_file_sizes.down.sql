--
-- DB migration file
--

BEGIN;

ALTER TABLE files
    DROP COLUMN file_size;

ALTER TABLE apt_files
    DROP COLUMN file_size;

COMMIT;

--
-- end of file
--