CREATE TABLE topicstoquotes(
   id SERIAL PRIMARY KEY,
   topicid int,
   quoteid int,
   created_at timestamptz default current_timestamp,
   updatedat timestamptz,
   deletedat timestamptz
);