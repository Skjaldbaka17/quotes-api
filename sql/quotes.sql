CREATE TABLE quotes(
   id SERIAL PRIMARY KEY,
   authorid integer not null,
   quote text NOT NULL,
   count integer default 0,
   isIcelandic boolean default false,
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz,
   FOREIGN KEY (authorid) REFERENCES authors(id) ON DELETE CASCADE
);