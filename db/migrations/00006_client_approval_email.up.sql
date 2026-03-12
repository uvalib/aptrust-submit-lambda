--
-- DB migration file
--

BEGIN;

ALTER TABLE clients
    ADD COLUMN approval_email VARCHAR( 64 ) NOT NULL DEFAULT '';

UPDATE clients set approval_email = 'dpg3k@virginia.edu' WHERE name = 'tracksys';

COMMIT;

--
-- end of file
--