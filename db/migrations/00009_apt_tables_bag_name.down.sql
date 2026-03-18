--
-- DB migration file
--

BEGIN;

ALTER TABLE apt_files
   ALTER COLUMN bag_name TYPE VARCHAR( 64 );

COMMIT;

--
-- end of file
--