ALTER TABLE "public"."events" ADD COLUMN "created_by" INTEGER NOT NULL REFERENCES "public"."users"("id") AFTER "end_at";
ALTER TABLE "public"."events" ADD COLUMN "updated_by" INTEGER NOT NULL REFERENCES "public"."users"("id") AFTER "created_by";
ALTER TABLE "public"."events" ADD COLUMN "deleted_by" INTEGER NOT NULL REFERENCES "public"."users"("id") AFTER "updated_by";
