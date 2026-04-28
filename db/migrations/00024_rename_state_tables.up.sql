--
-- DB migration file
--

BEGIN;

-- rename the types
ALTER TYPE submission_states RENAME TO submission_state_names;
ALTER TYPE bag_states RENAME TO bag_state_names;

-- rename the tables
ALTER TABLE submission_state RENAME TO submission_states;
ALTER TABLE bag_state RENAME TO bag_states;

COMMIT;

--
-- end of file
--