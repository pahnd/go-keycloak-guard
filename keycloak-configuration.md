# Cấu hình Keycloak

## 1. Cấu hình Realm

### 1.1. Tạo Realm mới
```json
{
  "realm": "my-realm",
  "enabled": true,
  "registrationAllowed": true,
  "resetPasswordAllowed": true,
  "verifyEmail": false,
  "loginTheme": "keycloak",
  "accessTokenLifespan": 300,
  "refreshTokenLifespan": 1800,
  "ssoSessionIdleTimeout": 1800,
  "ssoSessionMaxLifespan": 36000
}
```

### 1.2. Cấu hình Client
```json
{
  "clientId": "my-client",
  "enabled": true,
  "protocol": "openid-connect",
  "clientAuthenticatorType": "client-secret",
  "directAccessGrantsEnabled": true,
  "serviceAccountsEnabled": true,
  "authorizationServicesEnabled": true,
  "redirectUris": ["http://localhost:8000/*"],
  "webOrigins": ["http://localhost:8000"],
  "standardFlowEnabled": true,
  "implicitFlowEnabled": false,
  "clientSecret": "your-client-secret"
}
```

## 2. Cấu hình User

### 2.1. Tạo User
```json
{
  "username": "testuser",
  "enabled": true,
  "emailVerified": true,
  "firstName": "Test",
  "lastName": "User",
  "email": "test@example.com",
  "credentials": [
    {
      "type": "password",
      "value": "password123",
      "temporary": false
    }
  ],
  "requiredActions": [],
  "attributes": {
    "phone": ["+1234567890"],
    "department": ["IT"]
  }
}
```

### 2.2. Cấu hình Role
```json
{
  "name": "admin",
  "description": "Administrator role",
  "composite": false,
  "clientRole": false,
  "containerId": "my-realm",
  "attributes": {}
}
```

## 3. Cấu hình Authorization

### 3.1. Role-Based Access Control (RBAC)
```json
{
  "name": "resource-server",
  "type": "resource-server",
  "policies": [
    {
      "name": "admin-policy",
      "type": "role",
      "roles": ["admin"],
      "decisionStrategy": "AFFIRMATIVE"
    }
  ],
  "resources": [
    {
      "name": "api-resource",
      "type": "api",
      "owner": "resource-server",
      "uris": ["/api/*"]
    }
  ],
  "scopes": [
    {
      "name": "read",
      "displayName": "Read access"
    },
    {
      "name": "write",
      "displayName": "Write access"
    }
  ]
}
```

### 3.2. UMA (User-Managed Access)
```json
{
  "name": "uma-resource",
  "type": "uma",
  "owner": "resource-server",
  "resources": [
    {
      "name": "document",
      "type": "document",
      "owner": "resource-server",
      "scopes": ["read", "write", "delete"]
    }
  ],
  "policies": [
    {
      "name": "owner-policy",
      "type": "owner",
      "decisionStrategy": "AFFIRMATIVE"
    }
  ]
}
```

### 3.3. RPT (Requesting Party Token)
```json
{
  "name": "rpt-policy",
  "type": "rpt",
  "decisionStrategy": "AFFIRMATIVE",
  "permissions": [
    {
      "name": "resource-permission",
      "type": "resource",
      "resources": ["document"],
      "scopes": ["read", "write"]
    }
  ]
}
```

## 4. Cấu hình Token

### 4.1. JWT Token Configuration
```json
{
  "algorithm": "RS256",
  "keySize": 2048,
  "tokenLifespan": 300,
  "refreshTokenLifespan": 1800,
  "includeStandardClaims": true,
  "customClaims": {
    "department": "${user.attributes.department}",
    "phone": "${user.attributes.phone}"
  }
}
```

### 4.2. Token Exchange
```json
{
  "grantTypes": ["urn:ietf:params:oauth:grant-type:token-exchange"],
  "clientId": "token-exchange-client",
  "clientSecret": "exchange-secret",
  "directAccessGrantsEnabled": true,
  "serviceAccountsEnabled": true
}
```

## 5. Cấu hình Security

### 5.1. Password Policy
```json
{
  "passwordPolicy": "length(8) and requireNumbers and requireSpecialChars and requireUppercase and requireLowercase",
  "bruteForceProtected": true,
  "failureFactor": 5,
  "waitIncrementSeconds": 60,
  "quickLoginCheckMilliSeconds": 1000,
  "minimumQuickLoginWaitSeconds": 60,
  "maxFailureWaitSeconds": 900,
  "maxDeltaTimeSeconds": 43200,
  "failureResetTimeSeconds": 43200
}
```

### 5.2. SSL/TLS Configuration
```json
{
  "ssl": "external",
  "truststore": "/path/to/truststore.jks",
  "truststorePassword": "truststore-password",
  "keystore": "/path/to/keystore.jks",
  "keystorePassword": "keystore-password",
  "keyPassword": "key-password"
}
```

## 6. Cấu hình Federation

### 6.1. LDAP Integration
```json
{
  "name": "ldap-provider",
  "providerId": "ldap",
  "providerType": "org.keycloak.storage.UserStorageProvider",
  "parentId": "my-realm",
  "config": {
    "enabled": ["true"],
    "connectionUrl": ["ldap://ldap.example.com:389"],
    "bindDn": ["cn=admin,dc=example,dc=com"],
    "bindCredential": ["admin-password"],
    "baseDn": ["dc=example,dc=com"],
    "userDn": ["ou=users"],
    "searchScope": ["1"],
    "authType": ["simple"]
  }
}
```

### 6.2. OIDC Provider
```json
{
  "name": "oidc-provider",
  "providerId": "oidc",
  "providerType": "org.keycloak.broker.oidc.OIDCIdentityProvider",
  "parentId": "my-realm",
  "config": {
    "clientId": ["external-client"],
    "clientSecret": ["external-secret"],
    "authorizationUrl": ["https://external-idp.com/auth"],
    "tokenUrl": ["https://external-idp.com/token"],
    "userInfoUrl": ["https://external-idp.com/userinfo"],
    "defaultScopes": ["openid profile email"]
  }
}
```

## 7. Cấu hình Monitoring

### 7.1. Metrics Configuration
```json
{
  "metricsEnabled": true,
  "metricsPort": 9990,
  "metricsPath": "/metrics",
  "metricsRealm": "master",
  "metricsClientId": "metrics-client",
  "metricsClientSecret": "metrics-secret"
}
```

### 7.2. Health Check
```json
{
  "healthEnabled": true,
  "healthReadinessEnabled": true,
  "healthLivenessEnabled": true,
  "healthReadinessPath": "/health/ready",
  "healthLivenessPath": "/health/live"
}
``` 