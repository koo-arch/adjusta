-- Create "outbox_messages" table
CREATE TABLE "adjusta"."outbox_messages" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "message_type" character varying NOT NULL,
  "aggregate_type" character varying NOT NULL,
  "aggregate_id" uuid NOT NULL,
  "payload" jsonb NOT NULL DEFAULT '{}',
  "dispatch_attempts" bigint NOT NULL DEFAULT 0,
  "available_at" timestamptz NOT NULL,
  "dispatched_at" timestamptz NULL,
  "processed_at" timestamptz NULL,
  "last_dispatch_error" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "outboxmessage_available_at" to table: "outbox_messages"
CREATE INDEX "outboxmessage_available_at" ON "adjusta"."outbox_messages" ("available_at") WHERE (dispatched_at IS NULL AND processed_at IS NULL);
-- Create index "outboxmessage_processed_at" to table: "outbox_messages"
CREATE INDEX "outboxmessage_processed_at" ON "adjusta"."outbox_messages" ("processed_at");
-- Create index "outboxmessage_aggregate_type_aggregate_id_created_at" to table: "outbox_messages"
CREATE INDEX "outboxmessage_aggregate_type_aggregate_id_created_at" ON "adjusta"."outbox_messages" ("aggregate_type", "aggregate_id", "created_at");
