-- modify "queue_jobs" table
ALTER TABLE "queue_jobs" ADD COLUMN "status" character varying NOT NULL DEFAULT 'pending';
