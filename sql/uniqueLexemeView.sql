CREATE MATERIALIZED VIEW unique_lexeme AS
SELECT word FROM ts_stat('SELECT to_tsvector(''simple'', quotes.quote) || 
    to_tsvector(''simple'', authors.name) 
FROM quotes
JOIN authors ON authors.id = quotes.author_id
GROUP BY quotes.id, authors.id');