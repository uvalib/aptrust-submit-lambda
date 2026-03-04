--
-- DB migration file
--

BEGIN;

CREATE TYPE submission_states AS ENUM (
    'registered',               -- registered but not yet started
    'validating',               -- validating the submission
    'building',                 -- building the submission assets
    'pending-approval',         -- waiting for manual approval
    'submitting',               -- submitting to APTrust
    'pending-ingest',           -- waiting for APTrust to ingest
    'error',                    -- submission has errored in some manner
    'incomplete',               -- submission is incomplete (one or more bags rejected by APT)
    'complete');                -- submission is complete

CREATE TYPE bag_states AS ENUM (
    'registered',               -- registered but not yet started
    'building',                 -- building the bag
    'ready',                    -- built and ready to be submitted
    'submitting',               -- submitting to APTrust
    'pending-ingest',           -- waiting for APTrust to ingest
    'error',                    -- bag has errored in some manner
    'complete');                -- bag is complete

-- create clients table
CREATE TABLE clients (
    id         serial PRIMARY KEY,
    name       VARCHAR( 64 ) NOT NULL DEFAULT '',
    identifier VARCHAR( 64 ) NOT NULL DEFAULT '',
    created_at timestamp DEFAULT NOW()
);

-- create the clients index(s)
CREATE UNIQUE INDEX clients_identifier_distinct_idx ON clients(identifier);

-- create submissions table
CREATE TABLE submissions (
    id         serial PRIMARY KEY,
    identifier VARCHAR( 64 ) NOT NULL DEFAULT '',
    client     VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES clients(identifier),
    status     submission_states NOT NULL DEFAULT 'registered',
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW()
);

-- create the submissions index(s)
CREATE UNIQUE INDEX submissions_identifier_distinct_idx ON submissions(identifier);
CREATE INDEX submissions_client_key_idx ON submissions(client);

-- create bags table
CREATE TABLE bags (
    id            serial PRIMARY KEY,
    name          VARCHAR( 64 ) NOT NULL DEFAULT '', -- this is the external name
    submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
    status        bag_states NOT NULL DEFAULT 'registered',
    etag          VARCHAR( 64 ) NOT NULL DEFAULT '',
    created_at    timestamp DEFAULT NOW(),
    updated_at    timestamp DEFAULT NOW()
);

-- create the bags index(s)
--CREATE UNIQUE INDEX bags_identifier_distinct_idx ON bags(identifier);
CREATE INDEX bags_submission_key_idx ON bags(submission);

-- create files table
CREATE TABLE files (
    id            serial PRIMARY KEY,
    name          TEXT NOT NULL DEFAULT '',
    hash          VARCHAR( 64 ) NOT NULL DEFAULT '',
    submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
    bag_name      VARCHAR( 64 ) NOT NULL DEFAULT '',
    created_at    timestamp DEFAULT NOW()
);

-- create the files index(s)
CREATE INDEX files_submission_key_idx ON files(submission);
CREATE INDEX files_bag_name_key_idx ON files(bag_name);

-- create apt_files table
CREATE TABLE apt_files (
    id            serial PRIMARY KEY,
    file_name     TEXT NOT NULL DEFAULT '',
    hash          VARCHAR( 64 ) NOT NULL DEFAULT '',
    bag_name      VARCHAR( 64 ) NOT NULL DEFAULT '',
    apt_added_at  timestamp DEFAULT NOW(),
    created_at    timestamp DEFAULT NOW()
);

-- create the apt_files index(s)
CREATE INDEX apt_files_hash_key_idx ON apt_files(hash);
CREATE INDEX apt_files_hash_bag_name_key_idx ON apt_files(bag_name);

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

ALTER TABLE files SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE files SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE files SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE files SET (autovacuum_analyze_threshold = 1000);

ALTER TABLE apt_files SET (autovacuum_vacuum_scale_factor = 0.2);  -- 20%
ALTER TABLE apt_files SET (autovacuum_vacuum_threshold = 1000);
ALTER TABLE apt_files SET (autovacuum_analyze_scale_factor = 0.1); -- 10%
ALTER TABLE apt_files SET (autovacuum_analyze_threshold = 1000);

-- add the clients we know about
INSERT INTO clients(name, identifier) VALUES('avalon',     'cid-avalon-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('libra-data', 'cid-libra-data-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('libra-etd',  'cid-libra-etd-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('libra-open', 'cid-libra-open-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('tracksys',   'cid-tracksys-' || substr(md5(random()::text), 0, 11));

COMMIT;

--
-- end of file
--
