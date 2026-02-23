--
-- DB migration file
--

BEGIN;

-- drop the tables if they exist
DROP TABLE IF EXISTS bags;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS clients;

COMMIT;

--
-- end of file
--