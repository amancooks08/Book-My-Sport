CREATE TABLE "user"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "contact" VARCHAR(10) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "city" VARCHAR(255) NOT NULL,
    "state" VARCHAR(255) NOT NULL,
    "type" VARCHAR(255) NOT NULL
);


CREATE TABLE "venue"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "city" VARCHAR(255) NOT NULL,
    "state" VARCHAR(255) NOT NULL,
    "contact" VARCHAR(10) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "opening_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "closing_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "games" TEXT NOT NULL,
    "price" BIGINT NOT NULL,
    "rating" DECIMAL(8, 2) NOT NULL
);


CREATE TABLE "booking"(
    "id" SERIAL PRIMARY KEY,
    "booked_by" BIGINT NOT NULL REFERENCES "user"("id"),
    "booked_at" BIGINT NOT NULL REFERENCES "venue"("id"),
    "time" TIMESTAMP(0) WITH TIME zone NOT NULL,
    "start_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "end_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "duration" BIGINT NOT NULL,
    "game" TEXT NOT NULL,
    "amount" DOUBLE PRECISION NOT NULL
);


CREATE TABLE "slots"(
    "id" SERIAL PRIMARY KEY,
    "venue_id" BIGINT NOT NULL REFERENCES "venue"("id"),
    "start_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "end_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
    "duration" BIGINT NOT NULL,
    "status" VARCHAR(255) NOT NULL
);

-- ALTER TABLE
--     "venue" ADD CONSTRAINT "venue_bookings_foreign" FOREIGN KEY("bookings") REFERENCES "booking"("id");

-- ALTER TABLE
--     "booking" ADD CONSTRAINT "bookings_slot_id_foreign" FOREIGN KEY("slot_id") REFERENCES "slots"("id");

-- ALTER TABLE
--     "Venue" ADD CONSTRAINT "venue_slots_foreign" FOREIGN KEY("slots") REFERENCES "slots"("id");