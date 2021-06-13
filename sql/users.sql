CREATE TABLE users(
   id SERIAL PRIMARY KEY,
   email varchar not null unique,
   name VARCHAR NOT NULL,
   api_key_hash varchar not null unique,
   password_hash text not null,
   tier varchar not null default 'free',
   create_dat timestamptz,
   update_dat timestamptz,
   delete_dat timestamptz
);