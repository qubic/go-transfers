create table if not exists entities
(
    id         bigint primary key generated by default as identity,
    identity   text not null unique, -- indexed
    created_at timestamp with time zone default now() not null
);
