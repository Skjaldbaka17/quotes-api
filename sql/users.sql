CREATE TABLE users(
   id SERIAL PRIMARY KEY,
   email varchar not null unique,
   name VARCHAR NOT NULL,
   api_key varchar not null unique,
   password_hash text not null,
   tier varchar not null default 'free',
   created_at timestamptz default current_timestamp,
   updated_at timestamptz,
   deleted_at timestamptz
);