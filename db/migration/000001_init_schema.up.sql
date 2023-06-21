create table accounts(
                         id bigserial primary key,
                         owner varchar(100) not null,
                         balance bigint not null,
                         currency varchar(50) not null,
                         created_at timestamptz not null default now()
);

create table entries (
                         id bigserial primary key,
                         account_id bigint not null ,
                         amount bigint not null ,
                         created_at timestamptz not null default now()
);

create table transfers (
                           id bigserial primary key,
                           from_account_id bigint not null ,
                           to_account_id bigint not null ,
                           amount bigint not null check(amount > -1) ,
                           created_at timestamptz not null default now()
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

create index on accounts (owner);
create index on entries (account_id);
create index on transfers (from_account_id);
create index on transfers (to_account_id);
create index on transfers (from_account_id,to_account_id);
