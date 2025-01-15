create table if not exists schemaName.usr
(
    id              bigserial primary key,
    created         timestamp    default CURRENT_TIMESTAMP not null,
    modified        timestamp    default CURRENT_TIMESTAMP not null,
    email           varchar(255) unique,
    password        varchar(255)
);

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_trigger
            WHERE tgname = 'update_usr_modified_at_schemaName'
        ) THEN
            CREATE TRIGGER update_usr_modified_at_schemaName
                BEFORE UPDATE ON schemaName.usr
                FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
        END IF;
    END $$;



