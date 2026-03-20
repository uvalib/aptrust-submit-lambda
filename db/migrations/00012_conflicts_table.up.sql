--
-- DB migration file
--

BEGIN;

CREATE TYPE conflict_basis AS ENUM (
    'aptrust',           -- conflict with an aptrust file
    'local');            -- conflict with a local file

-- create conflicts table
CREATE TABLE conflicts (
   id               serial PRIMARY KEY,
   submission       VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
   new_file         INTEGER NOT NULL REFERENCES files(id),
   basis            conflict_basis NOT NULL DEFAULT 'aptrust',
   conflicting_file INTEGER NOT NULL,          -- could be either a files or an apt_files reference
   created_at       timestamp DEFAULT NOW()
);

-- create the conflicts index(s)
CREATE INDEX conflicts_submission_idx ON conflicts(submission);

COMMIT;

--
-- end of file
--