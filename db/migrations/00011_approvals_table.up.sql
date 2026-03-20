--
-- DB migration file
--

BEGIN;

-- create approvals table
CREATE TABLE approvals (
   id            serial PRIMARY KEY,
   submission    VARCHAR( 64 ) NOT NULL DEFAULT '' REFERENCES submissions(identifier),
   who           VARCHAR( 64 ) NOT NULL DEFAULT '',
   created_at    timestamp DEFAULT NOW()
);

-- create the approvals index(s)
CREATE INDEX approvals_submission_idx ON approvals(submission);

COMMIT;

--
-- end of file
--