CREATE TABLE IF NOT EXISTS schemaName.permissions (
                                                      id            bigserial PRIMARY KEY,
                                                      created       timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                                      modified      timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                                      name          varchar(255) UNIQUE NOT NULL,
    description   TEXT
    );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'update_permissions_modified_at_schemaName'
    ) THEN
CREATE TRIGGER update_permissions_modified_at_schemaName
    BEFORE UPDATE ON schemaName.permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
END IF;
END $$;
