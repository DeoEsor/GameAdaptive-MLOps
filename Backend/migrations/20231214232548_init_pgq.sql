-- +goose Up
-- +goose StatementBegin
CREATE TABLE pgq_jobs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  task_name TEXT NOT NULL,
  payload BYTEA NOT NULL,
  run_after TIMESTAMP WITH TIME ZONE NOT NULL,
  retry_waits INT NOT NULL,
  error TEXT,
  in_run BOOLEAN NOT NULL DEFAULT true
);

-- Add an index for fast fetching of jobs by queue_name, sorted by run_after.  But only
-- index jobs that haven't been done yet, in case the user is keeping the job history around.
CREATE INDEX idx_pgq_jobs_fetch
	ON pgq_jobs (run_after, retry_waits)
  WHERE in_run;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_pgq_jobs_fetch;
DROP TABLE pgq_jobs;
-- +goose StatementEnd
