CREATE TABLE IF NOT EXISTS schemaName.role_permissions
(
    role_id       bigint,
    permission_id bigint,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES schemaName.roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES schemaName.permissions(id) ON DELETE CASCADE
    );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'update_role_permissions_modified_at_schemaName'
    ) THEN
CREATE TRIGGER update_role_permissions_modified_at_schemaName
    BEFORE UPDATE ON schemaName.role_permissions
    FOR
    EACH ROW EXECUTE FUNCTION update_updated_at_column();
END IF;
END $$;
