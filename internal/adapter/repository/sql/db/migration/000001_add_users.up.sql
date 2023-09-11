CREATE TABLE "users" (
    "id" char(26) PRIMARY KEY,
    "email" varchar NOT NULL UNiQUE,
    "username" varchar NOT NULL UNiQUE,
    "password" varchar NOT NULL,
    "image" TEXT NOT NULL DEFAULT 'https://api.realworld.io/images/demo-avatar.png',
    "bio" TEXT NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

--bun:split
CREATE TABLE "user_follows" (
    "follower_id" char(26) NOT NULL,
    "followee_id" char(26) NOT NULL,
    PRIMARY KEY ("follower_id", "followee_id"),
    FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("followee_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)