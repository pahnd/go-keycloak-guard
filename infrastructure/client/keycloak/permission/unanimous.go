package permission

import (
	"strings"
)

type UnanimousDecisionStrategy struct {
}

func (d *UnanimousDecisionStrategy) HasPermissions(requiredPermissions []string, permissions PermissionCollection) bool {
	// Requires all permissions to be granted
	permissionsGranted := 0
	if len(requiredPermissions) == 0 {
		return true
	}
	for _, permission := range permissions {
		for _, scope := range permission.ResourceScopes {
			for _, requiredPermission := range requiredPermissions {
				requiredResource, requiredScope := strings.Split(requiredPermission, "#")[0], strings.Split(requiredPermission, "#")[1]
				if permission.ResourceName == requiredResource && scope == requiredScope {
					permissionsGranted++
				}
			}
		}
	}
	if len(requiredPermissions) > permissionsGranted {
		return false
	}
	return true
}
