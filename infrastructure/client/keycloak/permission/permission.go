package permission

type Permission struct {
	ResourceID     string   `json:"rsid"`
	ResourceName   string   `json:"rsname"`
	ResourceScopes []string `json:"scopes"`
}

type PermissionCollection []Permission

func (p *PermissionCollection) HasPermission(requiredPermissions []string, strategy string) (bool, error) {
	strategyClient := NewDecisionStrategyClient()
	return strategyClient.HasPermissions(requiredPermissions, *p, strategy)
}
