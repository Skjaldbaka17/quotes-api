CREATE TABLE qodice (
    id serial not null,
    quote_id integer not null,
    date date unique not null default current_date,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz
)