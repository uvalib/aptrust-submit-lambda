--
-- DB migration file
--

BEGIN;

-- drop the tables if they exist
DROP TABLE IF EXISTS bag_status;
DROP TABLE IF EXISTS submission_status;

-- drop the type definitions
DROP TYPE IF EXISTS bag_states;
DROP TYPE IF EXISTS submission_states;

COMMIT;

--
-- end of file
--