CREATE TABLE "articles" (
    "id" char(26) PRIMARY KEY,
    "slug" text NOT NULL,
    "title" text NOT NULL,
    "description" text NOT NULL,
    "body" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "author_id" char(26) NOT NULL,
    FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "comments" (
    "id" char(26) PRIMARY KEY,
    "body" TEXT NOT NULL,
    "article_id" char(26) NOT NULL,
    "author_id" char(26) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "tags" (
    "id" char(26) PRIMARY KEY,
    "name" varchar NOT NULL
);

--bun:split
CREATE TABLE "article_tags" (
    "article_id" char(26) NOT NULL,
    "tag_id" char(26) NOT NULL,
    PRIMARY KEY ("article_id", "tag_id"),
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

--bun:split
CREATE TABLE "article_favorites" (
    "user_id" char(26) NOT NULL,
    "article_id" char(26) NOT NULL,
    PRIMARY KEY ("user_id", "article_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
