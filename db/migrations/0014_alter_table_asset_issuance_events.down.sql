alter table asset_issuance_events
    alter column unit_of_measurement type bytea using decode(unit_of_measurement::text, 'base64');