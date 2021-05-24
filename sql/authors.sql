CREATE TABLE authors(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL,
   count integer default 0,
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz,
   tsv tsvector
);