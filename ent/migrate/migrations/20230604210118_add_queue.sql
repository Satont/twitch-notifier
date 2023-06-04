-- create "queue_jobs" table
CREATE TABLE "queue_jobs" ("id" uuid NOT NULL, "queue_name" character varying NOT NULL, "data" bytea NOT NULL, "retries" bigint NOT NULL, "max_retries" bigint NOT NULL, "added_at" timestamptz NOT NULL, "ttl" bigint NOT NULL, "fail_reason" character varying NULL, PRIMARY KEY ("id"));
-- create index "queuejob_queue_name" to table: "queue_jobs"
CREATE INDEX "queuejob_queue_name" ON "queue_jobs" ("queue_name");
