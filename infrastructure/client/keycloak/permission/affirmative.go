package permission

import (
	"strings"
)

type AffirmativeDecisionAStrategy struct {
}

func (d *AffirmativeDecisionAStrategy) HasPermissions(requiredPermissions []string, permissions PermissionCollection) bool {
	// Requires only one permission to be granted
	if len(requiredPermissions) == 0 {
		return true
	}
	for _, permission := range permissions {
		for _, scope := range permission.ResourceScopes {
			for _, requiredPermission := range requiredPermissions {
				requiredResource, requiredScope := strings.Split(requiredPermission, "#")[0], strings.Split(requiredPermission, "#")[1]
				if permission.ResourceName == requiredResource && scope == requiredScope {
					return true
				}
			}
		}
	}
	return false
}
