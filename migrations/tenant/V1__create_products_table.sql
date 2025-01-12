CREATE TABLE IF NOT EXISTS schemaName.products (
                                id SERIAL PRIMARY KEY,
                                created  timestamp default CURRENT_TIMESTAMP not null,
                                modified timestamp default CURRENT_TIMESTAMP not null,
                                name VARCHAR(255) NOT NULL,
                                description TEXT,
                                price DECIMAL(10, 2) NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_trigger
            WHERE tgname = 'update_products_modified_at_schemaName'
        ) THEN
            CREATE TRIGGER update_products_modified_at_schemaName
                BEFORE UPDATE ON schemaName.products
                FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
        END IF;
    END $$;
