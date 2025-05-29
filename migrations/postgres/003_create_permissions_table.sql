-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);

-- Create index on resource and action
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);

-- Create index on active status
CREATE INDEX IF NOT EXISTS idx_permissions_active ON permissions(is_active);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions(deleted_at);
