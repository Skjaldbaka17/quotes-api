UPDATE authors SET tsv = setweight(to_tsvector('english', name), 'A');
UPDATE quotes SET tsv = setweight(to_tsvector('english', quote), 'B');
CREATE INDEX if not exists index_authors_on_name ON authors USING gin(tsv);
CREATE INDEX if not exists index_quotes_on_quote ON quotes USING gin(tsv);
CREATE INDEX if not exists index_quotes_on_author_id ON quotes(author_id);
CREATE INDEX if not exists index_quotes_on_count ON quotes(count);

CREATE INDEX if not exists index_search_on_name_tsv ON searchview using gin(name_tsv);
CREATE INDEX if not exists index_search_on_quote_tsv ON searchview using gin(quote_tsv);
CREATE INDEX if not exists index_search_on_tsv ON searchview using gin(tsv);
CREATE INDEX if not exists index_search_on_author_id ON searchview(author_id);
CREATE INDEX if not exists index_search_on_quote_id ON searchview(quote_id);
CREATE INDEX if not exists index_search_on_quote_count ON searchview(quote_count);
CREATE INDEX if not exists index_search_on_author_count ON searchview(author_count);

CREATE INDEX if not exists index_topics_view_on_name_tsv ON topicsView using gin(name_tsv);
CREATE INDEX if not exists index_topics_view_on_quote_tsv ON topicsView using gin(quote_tsv);
CREATE INDEX if not exists index_topics_view_on_tsv ON topicsView using gin(tsv);
CREATE INDEX if not exists index_topics_view_on_author_id ON topicsView(author_id);
CREATE INDEX if not exists index_topics_view_on_quote_id ON topicsView(quote_id);

create INDEX if not exists index_request_history_on_user_id on requesthistory(user_id);
create INDEX if not exists index_request_history_on_created_at on requesthistory(created_at);

CREATE INDEX words_idx ON unique_lexeme USING gin(word gin_trgm_ops);
CREATE INDEX words_idx_quotes ON unique_lexeme_quotes USING gin(word gin_trgm_ops);
CREATE INDEX words_idx_authors ON unique_lexeme_authors USING gin(word gin_trgm_ops);

CREATE EXTENSION if not exists pg_trgm;

