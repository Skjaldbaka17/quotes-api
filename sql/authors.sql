CREATE TABLE authors(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL UNIQUE,
   count integer default 0,
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz,
   hasIcelandicQuotes boolean default false,
   tsv tsvector
);