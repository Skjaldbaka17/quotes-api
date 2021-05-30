create or replace view searchView as 
select authors.id as authorid,
       authors.name,
       quotes.id as quoteid,
       quotes.quote as quote,
       quotes.isicelandic as isicelandic,
       authors.tsv || quotes.tsv  as tsv,
       authors.tsv as nametsv,
       quotes.tsv as quotetsv
from authors
   inner join quotes
      on authors.id = quotes.authorid;