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
    'complete');                -- submission is complete

CREATE TYPE bag_states AS ENUM (
    'building',                 -- building the bag
    'ready',                    -- built and ready to be submitted
    'submitting',               -- submitting to APTrust
    'pending-ingest',           -- waiting for APTrust to ingest
    'error',                    -- bag has errored in some manner
    'complete');                -- bag is complete

-- create submission_status table
CREATE TABLE submission_status (
    id             serial PRIMARY KEY,
    submission     VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
    status         submission_states,
    created_at     timestamp DEFAULT NOW(),
    updated_at     timestamp DEFAULT NOW()
);

-- create bag_status table
CREATE TABLE bag_status (
    id          serial PRIMARY KEY,
    bag         VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES bags(identifier),
    status      bag_states,
    created_at  timestamp DEFAULT NOW(),
    updated_at  timestamp DEFAULT NOW()
);

-- create the index(s)
CREATE UNIQUE INDEX submission_status_distinct_idx ON submission_status(submission);
CREATE UNIQUE INDEX bag_status_distinct_idx ON bag_status(bag);

COMMIT;

--
-- end of file
--