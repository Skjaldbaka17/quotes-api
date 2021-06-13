CREATE TABLE errorhistory (
    id serial not null,
    user_id integer not null,
    route varchar not null,
    request_body text not null,
    error_message text not null,
    extra_info text,
    created_at timestamptz default current_timestamp
)