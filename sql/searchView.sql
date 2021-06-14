create or replace view searchView as 
select authors.id as author_id,
       authors.name,
       quotes.id as quote_id,
       quotes.quote as quote,
       quotes.is_icelandic as is_icelandic,
       authors.tsv || quotes.tsv  as tsv,
       authors.tsv as name_tsv,
       quotes.tsv as quote_tsv,
       quotes.count as quote_count,
       authors.count as author_count
from authors
   inner join quotes
      on authors.id = quotes.author_id;