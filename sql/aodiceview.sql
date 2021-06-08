create or replace view aodiceview as 
select a.id as authorid,
        a.name as name,
       aodice.date as date
from authors a
   inner join aodice
      on aodice.authorid = a.id;