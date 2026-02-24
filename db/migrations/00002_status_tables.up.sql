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
    submission_id  integer REFERENCES submissions(id),
    status         submission_states,
    created_at     timestamp DEFAULT NOW(),
    updated_at     timestamp DEFAULT NOW()
);

-- create bag_status table
CREATE TABLE bag_status (
    id          serial PRIMARY KEY,
    bag_id      integer REFERENCES bags(id),
    status      bag_states,
    created_at  timestamp DEFAULT NOW(),
    updated_at  timestamp DEFAULT NOW()
);

-- create the distinct indexes
CREATE UNIQUE INDEX submission_status_distinct_idx ON submission_status(submission_id);
CREATE UNIQUE INDEX bag_status_distinct_idx ON bag_status(bag_id);

COMMIT;

--
-- end of file
--