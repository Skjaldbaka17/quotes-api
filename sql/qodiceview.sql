create or replace view qodiceview as 
select q.id as quote_id,
        authors.name as name,
        q.quote as quote,
       authors.id as author_id,
       qodice.date as date
from authors
   inner join quotes q
      on authors.id = q.author_id
   inner join qodice
      on q.id = qodice.quote_id;