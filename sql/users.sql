CREATE TABLE users(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL UNIQUE,
   passwordhash text not null,
   email varchar,
   tier varchar not null default 'free',
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz
);