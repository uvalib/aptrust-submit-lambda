--
-- DB migration file
--

BEGIN;

ALTER TABLE apt_files
   ALTER COLUMN file_size TYPE INTEGER;

COMMIT;

--
-- end of file
--