-- add description column
ALTER TABLE schemaName.tenants ADD COLUMN IF NOT EXISTS  description TEXT;