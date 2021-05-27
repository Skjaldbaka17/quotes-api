UPDATE authors SET tsv = setweight(to_tsvector('english', name), 'A');
UPDATE quotes SET tsv = setweight(to_tsvector('english', quote), 'A');
CREATE INDEX if not exists index_authors_on_name ON authors USING gin(tsv);
CREATE INDEX if not exists index_quotes_on_quote ON quotes USING gin(tsv);
CREATE EXTENSION if not exists pg_trgm;
