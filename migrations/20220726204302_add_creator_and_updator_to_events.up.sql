ALTER TABLE "public"."events" ADD COLUMN "created_by" INTEGER NOT NULL REFERENCES "public"."users"("id");
ALTER TABLE "public"."events" ADD COLUMN "updated_by" INTEGER NOT NULL REFERENCES "public"."users"("id");
ALTER TABLE "public"."events" ADD COLUMN "deleted_by" INTEGER NULL REFERENCES "public"."users"("id");
