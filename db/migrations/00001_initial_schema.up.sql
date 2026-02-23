--
-- DB migration file
--

BEGIN;

-- create clients table
CREATE TABLE clients (
    id         serial PRIMARY KEY,
    name       VARCHAR( 64 ) NOT NULL DEFAULT '',
    identifier VARCHAR( 128 ) NOT NULL DEFAULT '',
    created_at timestamp DEFAULT NOW()
);

-- create submissions table
CREATE TABLE submissions (
    id         serial PRIMARY KEY,
    identifier VARCHAR( 128 ) NOT NULL DEFAULT '',
    client_id  integer REFERENCES clients(id),
    created_at timestamp DEFAULT NOW()
);

-- create bags table
CREATE TABLE bags (
    id            serial PRIMARY KEY,
    identifier    VARCHAR( 128 ) NOT NULL DEFAULT '',
    submission_id integer REFERENCES submissions(id),
    created_at    timestamp DEFAULT NOW()
);

-- create the submissions client_id key index
CREATE INDEX submissions_client_id_key_idx ON submissions(client_id);

-- create the bags submission_id key index
CREATE INDEX bags_submission_id_key_idx ON bags(submission_id);

-- auto vacuum parameters
-- see: https://aws.amazon.com/blogs/database/understanding-autovacuum-in-amazon-rds-for-postgresql-environments/

ALTER TABLE clients SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE clients SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE clients SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE clients SET (autovacuum_analyze_threshold = 1000);

ALTER TABLE submissions SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE submissions SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE submissions SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE submissions SET (autovacuum_analyze_threshold = 1000);

ALTER TABLE bags SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE bags SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE bags SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE bags SET (autovacuum_analyze_threshold = 1000);

-- add the clients we know about
INSERT INTO clients(name, identifier) VALUES('avalon',     'cid-' || gen_random_uuid());
INSERT INTO clients(name, identifier) VALUES('libra-etd',  'cid-' || gen_random_uuid());
INSERT INTO clients(name, identifier) VALUES('libra-open', 'cid-' || gen_random_uuid());
INSERT INTO clients(name, identifier) VALUES('tracksys',   'cid-' || gen_random_uuid());

COMMIT;

--
-- end of file
--
