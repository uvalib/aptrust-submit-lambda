--
-- DB migration file
--

BEGIN;

ALTER TABLE files
   ALTER COLUMN file_size TYPE BIGINT;

COMMIT;

--
-- end of file
--