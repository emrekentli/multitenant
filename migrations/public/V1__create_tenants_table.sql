CREATE TABLE IF NOT EXISTS schemaName.tenants (
                                id SERIAL PRIMARY KEY,
                                name VARCHAR(255) NOT NULL UNIQUE,
                                domain VARCHAR(255) NOT NULL UNIQUE,
                                schema_name VARCHAR(255) NOT NULL UNIQUE,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
