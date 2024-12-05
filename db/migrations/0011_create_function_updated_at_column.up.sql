create or replace function set_updated_at_time()
    returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language 'plpgsql';