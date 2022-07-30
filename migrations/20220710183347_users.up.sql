CREATE TABLE IF NOT EXISTS "public"."users" (
  "id" serial NOT NULL,
  PRIMARY KEY ("id"),
  "username" text NOT NULL,
  "password" text NULL,
  "display_name" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT current_timestamp,
  "updated_at" timestamp NOT NULL DEFAULT current_timestamp,
  "deleted_at" timestamp NULL
);

ALTER TABLE "public"."users" ADD CONSTRAINT "users_username" UNIQUE ("username");
