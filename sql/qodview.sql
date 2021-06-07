create or replace view qodview as 
select q.id as quoteid,
        authors.name as name,
        q.quote as quote,
       authors.id as authorid,
       qod.date as date
from authors
   inner join quotes q
      on authors.id = q.authorid
   inner join qod
      on q.id = qod.quoteid;