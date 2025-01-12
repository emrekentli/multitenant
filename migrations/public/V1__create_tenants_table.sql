create or replace function update_updated_at_column()
    returns trigger as $$
begin
    new.modified = now();
    return new;
end;
$$ language plpgsql;


CREATE TABLE IF NOT EXISTS schemaName.tenants (
                                id SERIAL PRIMARY KEY,
                                created  timestamp default CURRENT_TIMESTAMP not null,
                                modified timestamp default CURRENT_TIMESTAMP not null,
                                name VARCHAR(255) NOT NULL UNIQUE,
                                domain VARCHAR(255) NOT NULL UNIQUE,
                                schema_name VARCHAR(255) NOT NULL UNIQUE,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1
                       FROM pg_trigger
                       WHERE tgname = 'update_tenant_modified_at_schemaName') THEN
            create trigger update_tenant_modified_at_schemaName
                before update
                on schemaName.tenants
                for each row
            execute procedure update_updated_at_column();
        END IF;
    END
$$;
