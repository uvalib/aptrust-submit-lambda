--
-- DB migration file
--

BEGIN;

ALTER TABLE bags
   RENAME COLUMN name TO bag_name;

ALTER TABLE bag_state
    RENAME COLUMN name TO bag_name;

COMMIT;

--
-- end of file
--