CREATE TYPE "public"."field_type" AS ENUM (
  'text',
  'paragraph',
  'date',
  'number',
  'timestamp',
  'dropdown',
  'checkbox'
);

CREATE TABLE IF NOT EXISTS "public"."fields" (
  "id" serial NOT NULL,
  PRIMARY KEY ("id"),
  "workflow_id" integer NOT NULL REFERENCES "public"."workflows"("id"),
  "name" text NOT NULL,
  "type" field_type NOT NULL,
  "created_by" integer NOT NULL REFERENCES "public"."users"("id"),
  "updated_by" integer NOT NULL REFERENCES "public"."users"("id"),
  "deleted_by" integer NULL REFERENCES "public"."users"("id"),
  "created_at" timestamp NOT NULL DEFAULT current_timestamp,
  "updated_at" timestamp NOT NULL DEFAULT current_timestamp,
  "deleted_at" timestamp NULL
);
