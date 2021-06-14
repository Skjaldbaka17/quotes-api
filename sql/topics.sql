CREATE TABLE topics(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL UNIQUE,
   is_icelandic boolean default false,
   count integer default 0,
   created_at timestamptz default current_timestamp,
   updated_at timestamptz,
   deleted_at timestamptz,
   tsv tsvector
);