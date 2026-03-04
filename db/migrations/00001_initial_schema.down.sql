--
-- DB migration file
--

BEGIN;

-- drop the tables if they exist
DROP TABLE IF EXISTS apt_files;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS bags;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS clients;

-- drop the type definitions
DROP TYPE IF EXISTS bag_states;
DROP TYPE IF EXISTS submission_states;

COMMIT;

--
-- end of file
--