-- Add new schema named "adjusta"
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
CREATE SCHEMA "adjusta";
-- Create "accounts" table
CREATE TABLE "adjusta"."accounts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "google_user_id" character varying NOT NULL,
  "access_token" text NULL,
  "refresh_token" text NULL,
  "expires_at" timestamptz NULL,
  "scope" text NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "accounts_google_user_id_key" to table: "accounts"
CREATE UNIQUE INDEX "accounts_google_user_id_key" ON "adjusta"."accounts" ("google_user_id");
-- Create index "accounts_user_id_key" to table: "accounts"
CREATE UNIQUE INDEX "accounts_user_id_key" ON "adjusta"."accounts" ("user_id");
-- Create "calendars" table
CREATE TABLE "adjusta"."calendars" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "google_calendar_id" character varying NULL,
  "summary" character varying NULL,
  "description" text NULL,
  "timezone" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "calendars_google_calendar_id_key" to table: "calendars"
CREATE UNIQUE INDEX "calendars_google_calendar_id_key" ON "adjusta"."calendars" ("google_calendar_id");
-- Create "events" table
CREATE TABLE "adjusta"."events" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "title" character varying NOT NULL,
  "description" text NULL,
  "location" character varying NULL,
  "status" character varying NOT NULL DEFAULT 'active',
  "confirmed_google_event_id" character varying NULL,
  "sync_status" character varying NOT NULL DEFAULT 'not_synced',
  "last_synced_at" timestamptz NULL,
  "last_sync_error" text NULL,
  "primary_calendar_id" uuid NOT NULL,
  "confirmed_date_id" uuid NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "event_confirmed_date_id" to table: "events"
CREATE INDEX "event_confirmed_date_id" ON "adjusta"."events" ("confirmed_date_id");
-- Create index "event_primary_calendar_id" to table: "events"
CREATE INDEX "event_primary_calendar_id" ON "adjusta"."events" ("primary_calendar_id");
-- Create index "event_status" to table: "events"
CREATE INDEX "event_status" ON "adjusta"."events" ("status");
-- Create index "event_sync_status" to table: "events"
CREATE INDEX "event_sync_status" ON "adjusta"."events" ("sync_status");
-- Create index "event_user_id" to table: "events"
CREATE INDEX "event_user_id" ON "adjusta"."events" ("user_id");
-- Create "proposed_dates" table
CREATE TABLE "adjusta"."proposed_dates" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "google_event_id" character varying NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NOT NULL,
  "priority" bigint NOT NULL DEFAULT 0,
  "status" character varying NOT NULL DEFAULT 'active',
  "sync_status" character varying NOT NULL DEFAULT 'not_synced',
  "last_synced_at" timestamptz NULL,
  "last_sync_error" text NULL,
  "event_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "proposeddate_event_id" to table: "proposed_dates"
CREATE INDEX "proposeddate_event_id" ON "adjusta"."proposed_dates" ("event_id");
-- Create index "proposeddate_event_id_priority" to table: "proposed_dates"
CREATE UNIQUE INDEX "proposeddate_event_id_priority" ON "adjusta"."proposed_dates" ("event_id", "priority") WHERE (deleted_at IS NULL);
-- Create index "proposeddate_start_time" to table: "proposed_dates"
CREATE INDEX "proposeddate_start_time" ON "adjusta"."proposed_dates" ("start_time");
-- Create index "proposeddate_status" to table: "proposed_dates"
CREATE INDEX "proposeddate_status" ON "adjusta"."proposed_dates" ("status");
-- Create index "proposeddate_sync_status" to table: "proposed_dates"
CREATE INDEX "proposeddate_sync_status" ON "adjusta"."proposed_dates" ("sync_status");
-- Create "sessions" table
CREATE TABLE "adjusta"."sessions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "session_token" character varying NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "session_user_id" to table: "sessions"
CREATE INDEX "session_user_id" ON "adjusta"."sessions" ("user_id");
-- Create index "sessions_session_token_key" to table: "sessions"
CREATE UNIQUE INDEX "sessions_session_token_key" ON "adjusta"."sessions" ("session_token");
-- Create "user_calendars" table
CREATE TABLE "adjusta"."user_calendars" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "role" character varying NOT NULL,
  "is_visible" boolean NOT NULL DEFAULT true,
  "sync_proposed_dates" boolean NOT NULL DEFAULT false,
  "calendar_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "usercalendar_adjusta_candidate_user_id" to table: "user_calendars"
CREATE UNIQUE INDEX "usercalendar_adjusta_candidate_user_id" ON "adjusta"."user_calendars" ("user_id") WHERE (((role)::text = 'adjusta_candidate'::text) AND (deleted_at IS NULL));
-- Create index "usercalendar_calendar_id" to table: "user_calendars"
CREATE INDEX "usercalendar_calendar_id" ON "adjusta"."user_calendars" ("calendar_id");
-- Create index "usercalendar_primary_user_id" to table: "user_calendars"
CREATE UNIQUE INDEX "usercalendar_primary_user_id" ON "adjusta"."user_calendars" ("user_id") WHERE (((role)::text = 'primary'::text) AND (deleted_at IS NULL));
-- Create index "usercalendar_role" to table: "user_calendars"
CREATE INDEX "usercalendar_role" ON "adjusta"."user_calendars" ("role");
-- Create index "usercalendar_user_id" to table: "user_calendars"
CREATE INDEX "usercalendar_user_id" ON "adjusta"."user_calendars" ("user_id");
-- Create index "usercalendar_user_id_calendar_id" to table: "user_calendars"
CREATE UNIQUE INDEX "usercalendar_user_id_calendar_id" ON "adjusta"."user_calendars" ("user_id", "calendar_id");
-- Create "users" table
CREATE TABLE "adjusta"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "email" character varying NOT NULL,
  "name" character varying NULL,
  "avatar_url" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "adjusta"."users" ("email");
-- Modify "accounts" table
ALTER TABLE "adjusta"."accounts" ADD CONSTRAINT "accounts_users_account" FOREIGN KEY ("user_id") REFERENCES "adjusta"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "events" table
ALTER TABLE "adjusta"."events" ADD CONSTRAINT "events_calendars_primary_events" FOREIGN KEY ("primary_calendar_id") REFERENCES "adjusta"."calendars" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "events_proposed_dates_confirmed_date" FOREIGN KEY ("confirmed_date_id") REFERENCES "adjusta"."proposed_dates" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "events_users_events" FOREIGN KEY ("user_id") REFERENCES "adjusta"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "proposed_dates" table
ALTER TABLE "adjusta"."proposed_dates" ADD CONSTRAINT "proposed_dates_events_proposed_dates" FOREIGN KEY ("event_id") REFERENCES "adjusta"."events" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "sessions" table
ALTER TABLE "adjusta"."sessions" ADD CONSTRAINT "sessions_users_sessions" FOREIGN KEY ("user_id") REFERENCES "adjusta"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "user_calendars" table
ALTER TABLE "adjusta"."user_calendars" ADD CONSTRAINT "user_calendars_calendars_user_calendars" FOREIGN KEY ("calendar_id") REFERENCES "adjusta"."calendars" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "user_calendars_users_user_calendars" FOREIGN KEY ("user_id") REFERENCES "adjusta"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
