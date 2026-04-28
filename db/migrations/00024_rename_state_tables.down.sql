--
-- DB migration file
--

BEGIN;

-- rename the tables
ALTER TABLE bag_states RENAME TO bag_state;
ALTER TABLE submission_states RENAME TO submission_state;

-- rename the types
ALTER TYPE submission_state_names RENAME TO submission_states;
ALTER TYPE bag_state_names RENAME TO bag_states;

COMMIT;

--
-- end of file
--