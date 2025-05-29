-- Insert default roles
INSERT INTO roles (name, description, is_active) VALUES 
    ('admin', 'System Administrator with full access', true),
    ('hr_manager', 'HR Manager with employee management access', true),
    ('employee', 'Regular employee with basic access', true),
    ('viewer', 'Read-only access to basic information', true)
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (name, description, resource, action, is_active) VALUES 
    -- User management permissions
    ('users.create', 'Create new users', 'users', 'create', true),
    ('users.read', 'View user information', 'users', 'read', true),
    ('users.update', 'Update user information', 'users', 'update', true),
    ('users.delete', 'Delete users', 'users', 'delete', true),
    ('users.list', 'List all users', 'users', 'list', true),
    
    -- Role management permissions
    ('roles.create', 'Create new roles', 'roles', 'create', true),
    ('roles.read', 'View role information', 'roles', 'read', true),
    ('roles.update', 'Update role information', 'roles', 'update', true),
    ('roles.delete', 'Delete roles', 'roles', 'delete', true),
    ('roles.list', 'List all roles', 'roles', 'list', true),
    ('roles.assign', 'Assign roles to users', 'roles', 'assign', true),
    
    -- Permission management permissions
    ('permissions.create', 'Create new permissions', 'permissions', 'create', true),
    ('permissions.read', 'View permission information', 'permissions', 'read', true),
    ('permissions.update', 'Update permission information', 'permissions', 'update', true),
    ('permissions.delete', 'Delete permissions', 'permissions', 'delete', true),
    ('permissions.list', 'List all permissions', 'permissions', 'list', true),
    ('permissions.assign', 'Assign permissions to roles', 'permissions', 'assign', true),
    
    -- Profile permissions
    ('profile.read', 'View own profile', 'profile', 'read', true),
    ('profile.update', 'Update own profile', 'profile', 'update', true),
    
    -- System permissions
    ('system.admin', 'Full system administration access', 'system', 'admin', true),
    ('system.reports', 'Access to system reports', 'system', 'reports', true)
ON CONFLICT (name) DO NOTHING;
