create or replace view aodview as 
select a.id as id,
        a.name as name,
       aod.date as date
from authors a
   inner join aod
      on aod.authorid = a.id;