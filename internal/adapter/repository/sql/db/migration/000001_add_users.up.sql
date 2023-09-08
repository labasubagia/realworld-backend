CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "email" varchar NOT NULL UNiQUE,
    "username" varchar NOT NULL UNiQUE,
    "password" varchar NOT NULL,
    "image" TEXT NOT NULL DEFAULT 'https://realworld-temp-api.herokuapp.com/images/smiley-cyrus.jpeg',
    "bio" TEXT NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);