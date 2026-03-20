--
-- DB migration file
--

BEGIN;

-- drop the table if it exists
DROP TABLE IF EXISTS conflicts;

-- drop the type definitions
DROP TYPE IF EXISTS conflict_basis;

COMMIT;

--
-- end of file
--