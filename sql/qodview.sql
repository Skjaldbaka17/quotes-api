create or replace view qodview as 
select q.id as quote_id,
        authors.name as name,
        q.quote as quote,
       authors.id as author_id,
       qod.date as date,
       q.is_icelandic as is_icelandic
from authors
   inner join quotes q
      on authors.id = q.author_id
   inner join qod
      on q.id = qod.quote_id;