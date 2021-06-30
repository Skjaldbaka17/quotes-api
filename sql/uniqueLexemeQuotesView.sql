CREATE MATERIALIZED VIEW unique_lexeme_quotes AS
SELECT word FROM ts_stat('SELECT to_tsvector(''simple'', quotes.quote)
FROM quotes');