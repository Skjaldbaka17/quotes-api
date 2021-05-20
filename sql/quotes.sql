CREATE TABLE quotes(
   id SERIAL PRIMARY KEY,
   author_id integer not null,
   quote text NOT NULL,
   count integer default 0,
   created_at timestamptz,
   FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
);