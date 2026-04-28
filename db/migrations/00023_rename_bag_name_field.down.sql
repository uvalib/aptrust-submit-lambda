--
-- DB migration file
--

BEGIN;

ALTER TABLE bags
   RENAME COLUMN bag_name TO name;

ALTER TABLE bag_state
    RENAME COLUMN bag_name TO name;

COMMIT;

--
-- end of file
--