CREATE TABLE "sessions"
(
    "id"            uuid primary key,
    "username"      varchar     not null,
    "refresh_token" varchar     NOT NULL,
    "user_agent"    varchar     NOT NULL,
    "client_ip"     varchar     NOT NULL,
    "is_blocked"    boolean     not null default false,
    "expires_at"    timestamptz NOT NULL,
    "created_at"    timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("username") REFERENCES "users" ("username");