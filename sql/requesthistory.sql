CREATE TABLE requesthistory (
    id serial not null,
    user_id integer not null,
    api_key varchar not null,
    route varchar not null,
    request_body text not null,
    request text not null,
    created_at timestamptz default current_timestamp
)