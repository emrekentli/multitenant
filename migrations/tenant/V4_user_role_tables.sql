CREATE TABLE IF NOT EXISTS schemaName.user_roles (
                                                     user_id       bigint NOT NULL,
                                                     role_id       bigint NOT NULL,
                                                     created       timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                                     PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES schemaName.usr(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES schemaName.roles(id) ON DELETE CASCADE
    );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'update_user_roles_modified_at_schemaName'
    ) THEN
CREATE TRIGGER update_user_roles_modified_at_schemaName
    BEFORE UPDATE ON schemaName.user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
END IF;
END $$;
