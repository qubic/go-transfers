create table if not exists asset_change_events (
    id bigint primary key generated by default as identity,
    event_id bigint references events(id) unique not null,
    asset_id bigint references assets(id) not null,
    source_entity_id bigint references entities(id) not null,
    destination_entity_id bigint references entities(id) not null,
    number_of_shares bigint not null
);

create index on asset_change_events(event_id);
create index on asset_change_events(asset_id);
create index on asset_change_events(source_entity_id);
create index on asset_change_events(destination_entity_id);
