CREATE TABLE aod (
    id serial not null,
    authorid integer not null,
    date date unique not null default current_date,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz
)