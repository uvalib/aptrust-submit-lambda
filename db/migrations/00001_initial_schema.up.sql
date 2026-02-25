--
-- DB migration file
--

BEGIN;

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
    created_at timestamp DEFAULT NOW()
);

-- create the submissions index(s)
CREATE UNIQUE INDEX submissions_identifier_distinct_idx ON submissions(identifier);
CREATE INDEX submissions_client_key_idx ON submissions(client);

-- create bags table
CREATE TABLE bags (
    id            serial PRIMARY KEY,
    name          VARCHAR( 64 ) NOT NULL DEFAULT '', -- this is the external name
    identifier    VARCHAR( 64 ) NOT NULL DEFAULT '',
    submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
    created_at    timestamp DEFAULT NOW()
);

-- create the bags index(s)
CREATE UNIQUE INDEX bags_identifier_distinct_idx ON bags(identifier);
CREATE INDEX bags_submission_key_idx ON bags(submission);

-- create files table
CREATE TABLE files (
    id            serial PRIMARY KEY,
    name          TEXT NOT NULL DEFAULT '',
    hash          VARCHAR( 64 ) NOT NULL DEFAULT '',
    submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
    bag           VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES bags(identifier),
    created_at    timestamp DEFAULT NOW()
);

-- create the files index(s)
CREATE INDEX files_submission_key_idx ON files(submission);
CREATE INDEX files_bag_key_idx ON files(bag);

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

-- add the clients we know about
INSERT INTO clients(name, identifier) VALUES('avalon',     'cid-avalon-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('libra-etd',  'cid-libra-etd-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('libra-open', 'cid-libra-open-' || substr(md5(random()::text), 0, 11));
INSERT INTO clients(name, identifier) VALUES('tracksys',   'cid-tracksys-' || substr(md5(random()::text), 0, 11));

COMMIT;

--
-- end of file
--
