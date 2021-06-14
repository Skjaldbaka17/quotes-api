CREATE TABLE authors(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL UNIQUE,
   count integer default 0,
   created_at timestamptz default current_timestamp,
   updated_at timestamptz,
   deleted_at timestamptz,
   nr_of_english_quotes integer default 0,
   nr_of_icelandic_quotes integer default 0,
   has_icelandic_quotes boolean default false,
   tsv tsvector
);