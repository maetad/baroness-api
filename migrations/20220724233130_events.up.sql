CREATE TYPE "public"."event_platform" AS ENUM ('online', 'on_ground');
CREATE TYPE "public"."event_channel" AS ENUM ('record', 'live');

CREATE TABLE IF NOT EXISTS "public"."events" (
  "id" serial NOT NULL,
  PRIMARY KEY ("id"),
  "name" text NOT NULL,
  "platform" event_platform[],
  "channel" event_channel[],
  "start_at" timestamp NULL,
  "end_at" timestamp NULL,
  "created_at" timestamp NOT NULL DEFAULT current_timestamp,
  "updated_at" timestamp NOT NULL DEFAULT current_timestamp,
  "deleted_at" timestamp NULL
);
