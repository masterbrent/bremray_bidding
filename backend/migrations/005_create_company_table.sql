-- Create company table
CREATE TABLE IF NOT EXISTS companies (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo TEXT,
    address VARCHAR(255) DEFAULT '',
    city VARCHAR(100) DEFAULT '',
    state VARCHAR(2) DEFAULT '',
    zip VARCHAR(10) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    email VARCHAR(255) NOT NULL,
    license VARCHAR(100) DEFAULT '',
    website VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update updated_at
CREATE OR REPLACE FUNCTION update_companies_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER companies_updated_at_trigger
BEFORE UPDATE ON companies
FOR EACH ROW
EXECUTE FUNCTION update_companies_updated_at();

-- Insert default company
INSERT INTO companies (id, name, email) 
VALUES ('default', 'Bremray Electrical', 'info@bremray.com')
ON CONFLICT (id) DO NOTHING;