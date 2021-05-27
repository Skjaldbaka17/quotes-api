CREATE TABLE topicstoquotes(
   id SERIAL PRIMARY KEY,
   topicid int,
   quoteid int,
   createdat timestamptz,
   updatedat timestamptz,
   deletedat timestamptz
);