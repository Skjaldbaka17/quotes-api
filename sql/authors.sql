CREATE TABLE authors(
   id SERIAL PRIMARY KEY,
   name VARCHAR NOT NULL,
   count integer default 0,
   created_at timestamptz
);