--
-- DB migration file
--

BEGIN;

ALTER TABLE files
    ADD COLUMN file_size INTEGER NOT NULL DEFAULT 0;

ALTER TABLE apt_files
    ADD COLUMN file_size INTEGER NOT NULL DEFAULT 0;

COMMIT;

--
-- end of file
--