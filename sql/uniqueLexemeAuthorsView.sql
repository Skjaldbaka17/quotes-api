CREATE MATERIALIZED VIEW unique_lexeme_authors AS
SELECT word FROM ts_stat('SELECT to_tsvector(''simple'', authors.name)
FROM authors');