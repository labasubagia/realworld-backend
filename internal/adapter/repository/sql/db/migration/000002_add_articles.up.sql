CREATE TABLE "articles" (
    "id" serial PRIMARY KEY,
    "slug" text NOT NULL,
    "title" text NOT NULL,
    "description" text NOT NULL,
    "body" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "author_id" INTEGER NOT NULL,
    FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "comments" (
    "id" SERIAL PRIMARY KEY,
    "body" TEXT NOT NULL,
    "article_id" INTEGER NOT NULL,
    "author_id" INTEGER NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "tags" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar NOT NULL
);

--bun:split
CREATE TABLE "article_tags" (
    "article_id" INTEGER NOT NULL,
    "tag_id" INTEGER NOT NULL,
    PRIMARY KEY ("article_id", "tag_id"),
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "article_favorites" (
    "user_id" INTEGER NOT NULL,
    "article_id" INTEGER NOT NULL,
    PRIMARY KEY ("user_id", "article_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
