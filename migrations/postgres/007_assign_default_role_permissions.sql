-- Assign permissions to default roles

-- Admin role gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- HR Manager role gets user and basic management permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'hr_manager'
AND p.name IN (
    'users.create', 'users.read', 'users.update', 'users.list',
    'roles.read', 'roles.list', 'roles.assign',
    'permissions.read', 'permissions.list',
    'profile.read', 'profile.update',
    'system.reports'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Employee role gets basic permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'employee'
AND p.name IN (
    'users.read',
    'profile.read', 'profile.update'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Viewer role gets read-only permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'viewer'
AND p.name IN (
    'profile.read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
