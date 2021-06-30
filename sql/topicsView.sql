CREATE MATERIALIZED VIEW topicsView as 
select authors.id as author_id,
       authors.name,
       q.id as quote_id,
       q.quote as quote,
       q.is_icelandic as is_icelandic,
       authors.tsv || q.tsv  as tsv,
       authors.tsv as name_tsv,
       q.tsv as quote_tsv,
       t.name as topic_name,
       t.id as topic_id
from authors
   inner join quotes q
      on authors.id = q.author_id
   inner join topicstoquotes ttq
      on q.id = ttq.quote_id
   inner join topics t
      on t.id = ttq.topic_id;