# Script Workflow Xác thực với Keycloak

## 1. Role-Based Access Control (RBAC)

### 1.1. Cấu hình Realm và Client
```bash
# Tạo Realm mới
curl -X POST "http://localhost:8080/admin/realms" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "realm": "my-realm",
    "enabled": true
  }'

# Tạo Client
curl -X POST "http://localhost:8080/admin/realms/my-realm/clients" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "my-client",
    "enabled": true,
    "protocol": "openid-connect",
    "clientAuthenticatorType": "client-secret",
    "directAccessGrantsEnabled": true,
    "serviceAccountsEnabled": true,
    "authorizationServicesEnabled": true,
    "redirectUris": ["http://localhost:8000/*"],
    "webOrigins": ["http://localhost:8000"]
  }'
```

### 1.2. Tạo Role và Gán cho User
```bash
# Tạo Role
curl -X POST "http://localhost:8080/admin/realms/my-realm/roles" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "admin",
    "description": "Administrator role",
    "composite": false
  }'

# Gán Role cho User
curl -X POST "http://localhost:8080/admin/realms/my-realm/users/$USER_ID/role-mappings/realm" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{
    "id": "$ROLE_ID",
    "name": "admin",
    "composite": false,
    "clientRole": false,
    "containerId": "my-realm"
  }]'
```

## 2. User-Managed Access (UMA)

### 2.1. Cấu hình Resource Server
```bash
# Tạo Resource Server
curl -X POST "http://localhost:8080/realms/my-realm/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=my-client" \
  -d "client_secret=$CLIENT_SECRET"

# Tạo Resource
curl -X POST "http://localhost:8080/realms/my-realm/authz/protection/resource_set" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "document",
    "type": "document",
    "owner": "$USER_ID",
    "scopes": ["read", "write", "delete"]
  }'
```

### 2.2. Cấu hình Permission
```bash
# Tạo Permission
curl -X POST "http://localhost:8080/realms/my-realm/authz/protection/permission" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "document-permission",
    "type": "resource",
    "resources": ["document"],
    "scopes": ["read", "write"],
    "policies": ["owner-policy"]
  }'
```

## 3. Requesting Party Token (RPT)

### 3.1. Lấy RPT Token
```bash
# Lấy Permission Ticket
curl -X POST "http://localhost:8080/realms/my-realm/authz/protection/permission" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "$RESOURCE_ID",
    "resource_scopes": ["read", "write"]
  }'

# Đổi Permission Ticket lấy RPT
curl -X POST "http://localhost:8080/realms/my-realm/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=urn:ietf:params:oauth:grant-type:uma-ticket" \
  -d "ticket=$PERMISSION_TICKET" \
  -d "client_id=my-client" \
  -d "client_secret=$CLIENT_SECRET"
```

### 3.2. Sử dụng RPT Token
```bash
# Gọi API với RPT Token
curl -X GET "http://localhost:8000/api/resource" \
  -H "Authorization: Bearer $RPT_TOKEN"
```

## 4. Kết hợp các phương thức

### 4.1. Cấu hình Aggregate Policy
```bash
# Tạo Aggregate Policy
curl -X POST "http://localhost:8080/realms/my-realm/authz/protection/policy/aggregate" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "combined-policy",
    "policies": ["role-policy", "uma-policy", "rpt-policy"],
    "decisionStrategy": "AFFIRMATIVE"
  }'
```

### 4.2. Kiểm tra quyền truy cập
```bash
# Kiểm tra quyền với Aggregate Policy
curl -X POST "http://localhost:8080/realms/my-realm/authz/protection/policy/evaluate" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resources": ["document"],
    "context": {
      "attributes": {
        "user_id": "$USER_ID",
        "roles": ["admin"]
      }
    }
  }'
```

## 5. Script Kiểm tra và Monitoring

### 5.1. Kiểm tra Token
```bash
# Kiểm tra token hợp lệ
curl -X POST "http://localhost:8080/realms/my-realm/protocol/openid-connect/token/introspect" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=$TOKEN" \
  -d "client_id=my-client" \
  -d "client_secret=$CLIENT_SECRET"
```

### 5.2. Monitoring
```bash
# Kiểm tra health
curl -X GET "http://localhost:8080/health/ready"

# Lấy metrics
curl -X GET "http://localhost:8080/metrics" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 6. Script Backup và Recovery

### 6.1. Backup Configuration
```bash
# Export Realm
curl -X GET "http://localhost:8080/admin/realms/my-realm/export" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -o "realm-backup.json"

# Backup Users
curl -X GET "http://localhost:8080/admin/realms/my-realm/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -o "users-backup.json"
```

### 6.2. Recovery
```bash
# Import Realm
curl -X POST "http://localhost:8080/admin/realms" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d @realm-backup.json

# Restore Users
curl -X POST "http://localhost:8080/admin/realms/my-realm/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d @users-backup.json
```

## Lưu ý:
1. Thay thế các biến môi trường ($ADMIN_TOKEN, $CLIENT_SECRET, etc.) bằng giá trị thực
2. Đảm bảo Keycloak server đang chạy và có thể truy cập
3. Kiểm tra quyền truy cập trước khi thực hiện các lệnh
4. Backup dữ liệu trước khi thực hiện các thay đổi
5. Test các script trong môi trường development trước khi sử dụng trong production 