CREATE TABLE topics(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL UNIQUE,
   isIcelandic boolean default false,
   count integer default 0,
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz,
   tsv tsvector
);