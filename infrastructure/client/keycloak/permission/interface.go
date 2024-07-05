package permission

type DecisionInterface interface {
	HasPermissions(required []string, requiredPermissions PermissionCollection) bool
}
