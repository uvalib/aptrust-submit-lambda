--
-- DB migration file
--

BEGIN;

ALTER TABLE apt_files
   ALTER COLUMN bag_name TYPE TEXT;

COMMIT;

--
-- end of file
--