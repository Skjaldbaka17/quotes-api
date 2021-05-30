create or replace view TopicsView as 
select authors.id as authorid,
       authors.name,
       q.id as quoteid,
       q.quote as quote,
       q.isicelandic as isicelandic,
       authors.tsv || q.tsv  as tsv,
       authors.tsv as nametsv,
       q.tsv as quotetsv,
       t.name as topicname,
       t.id as topicid
from authors
   inner join quotes q
      on authors.id = q.authorid
   inner join topicstoquotes ttq
      on q.id = ttq.quoteid
   inner join topics t
      on t.id = ttq.topicid;