-- Table Definition
CREATE TABLE IF NOT EXISTS users (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "hashed_password" varchar(255) NOT NULL,
    "role" varchar(6) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);