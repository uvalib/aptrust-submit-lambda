--
-- DB migration file
--

BEGIN;

ALTER TABLE clients
    ADD CONSTRAINT clients_storage_fkey
    FOREIGN KEY (default_storage)
    REFERENCES storage_options(value);

ALTER TABLE submissions
    ADD CONSTRAINT submissions_storage_fkey
    FOREIGN KEY (storage)
    REFERENCES storage_options(value);

COMMIT;

--
-- end of file
--