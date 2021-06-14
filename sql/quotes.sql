CREATE TABLE quotes(
   id SERIAL PRIMARY KEY,
   author_id integer not null,
   quote text NOT NULL unique,
   count integer default 0,
   is_icelandic boolean default false,
   created_at timestamptz default current_timestamp,
   updated_at timestamptz,
   deleted_at timestamptz,
   tsv tsvector,
   FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
);