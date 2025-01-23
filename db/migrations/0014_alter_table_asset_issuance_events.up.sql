alter table asset_issuance_events
    alter column unit_of_measurement type text using encode(unit_of_measurement::bytea, 'base64');