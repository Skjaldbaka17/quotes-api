create or replace view qodiceview as 
select q.id as quoteid,
        authors.name as name,
        q.quote as quote,
       authors.id as authorid,
       qodice.date as date
from authors
   inner join quotes q
      on authors.id = q.authorid
   inner join qodice
      on q.id = qodice.quoteid;