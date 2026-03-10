--
-- DB migration file
--

BEGIN;

ALTER TABLE clients
    ADD COLUMN default_storage VARCHAR( 64 ) NOT NULL DEFAULT 'Standard';

ALTER TABLE submissions
    ADD COLUMN storage VARCHAR( 64 ) NOT NULL DEFAULT 'Standard';

-- set the default for tracksys
UPDATE clients set default_storage = 'Glacier-Deep-OH' WHERE name = 'tracksys';

COMMIT;

--
-- end of file
--