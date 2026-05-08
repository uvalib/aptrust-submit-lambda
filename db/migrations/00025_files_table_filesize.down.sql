--
-- DB migration file
--

BEGIN;

ALTER TABLE files
   ALTER COLUMN file_size TYPE INTEGER;

COMMIT;

--
-- end of file
--