--
-- DB migration file
--

BEGIN;

ALTER TYPE submission_states ADD VALUE 'abandoned';
ALTER TYPE bag_states ADD VALUE 'abandoned';

COMMIT;

--
-- end of file
--