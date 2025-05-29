# rbac/ - Control de Acceso Basado en Roles

Implementación de RBAC (Role-Based Access Control) para gestión granular de permisos.

## Responsabilidades

- Definir y gestionar roles
- Asignar permisos a roles
- Verificar permisos de usuario
- Implementar políticas de acceso
- Integrar con sistemas de autorización

## Archivos

- **`enforcer.go`** - Verificador principal de permisos
- **`policy.go`** - Definición de políticas de acceso
- **`adapter.go`** - Adaptador para persistencia (Casbin)
- **`service.go`** - Servicio de RBAC

## Modelo de Permisos

### Recursos
- `employee` - Gestión de empleados
- `user` - Gestión de usuarios
- `department` - Gestión de departamentos
- `report` - Acceso a reportes

### Acciones
- `create` - Crear recursos
- `read` - Leer/consultar recursos
- `update` - Modificar recursos
- `delete` - Eliminar recursos

### Roles Predefinidos
- `admin` - Acceso completo al sistema
- `hr_manager` - Gestión de RH
- `department_manager` - Gestión departamental
- `employee` - Acceso básico

## Formato de Permisos

```
{role}, {resource}:{action}
```

Ejemplos:
- `hr_manager, employee:create`
- `department_manager, employee:read`
- `admin, *:*` (acceso total)

## Implementación

```go
type RBACService struct {
    enforcer *casbin.Enforcer
}

func (r *RBACService) HasPermission(
    userRoles []string, 
    resource, action string,
) bool {
    permission := fmt.Sprintf("%s:%s", resource, action)
    
    for _, role := range userRoles {
        if r.enforcer.Enforce(role, permission) {
            return true
        }
    }
    return false
}

func (r *RBACService) AssignRole(userID uuid.UUID, role string) error
func (r *RBACService) RemoveRole(userID uuid.UUID, role string) error
func (r *RBACService) GetUserRoles(userID uuid.UUID) ([]string, error)
```

## Configuración

Archivo de modelo Casbin (`rbac_model.conf`):
```ini
[request_definition]
r = sub, obj

[policy_definition]
p = sub, obj

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj
```

## Uso

```go
// Verificar permiso en middleware
if !rbacService.HasPermission(userRoles, "employee", "create") {
    return fiber.ErrForbidden
}

// Asignar rol
rbacService.AssignRole(userID, "hr_manager")
```
