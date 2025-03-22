# Hướng dẫn về Authorization trong Keycloak

## 1. Role-Based Access Control (RBAC)

### 1.1. Cấu hình cơ bản
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
  ]
}
```

### 1.2. Tác dụng đối với User
- **Phân quyền theo Role**: User được gán các role cụ thể (ví dụ: admin, user, manager)
- **Kiểm soát truy cập**: Chỉ những user có role phù hợp mới được phép truy cập tài nguyên
- **Dễ quản lý**: Có thể gán/bỏ role cho user một cách linh hoạt
- **Phân cấp rõ ràng**: Mỗi role có quyền hạn riêng biệt

### 1.3. Ví dụ thực tế
- Admin có quyền truy cập tất cả tài nguyên
- User thường chỉ có quyền đọc
- Manager có quyền quản lý nhóm user

## 2. User-Managed Access (UMA)

### 2.1. Cấu hình cơ bản
```json
{
  "name": "uma-resource",
  "type": "uma",
  "resources": [
    {
      "name": "document",
      "type": "document",
      "scopes": ["read", "write", "delete"]
    }
  ]
}
```

### 2.2. Tác dụng đối với User
- **Quản lý quyền truy cập**: User có thể tự quản lý quyền truy cập tài nguyên của mình
- **Chia sẻ tài nguyên**: User có thể chia sẻ tài nguyên với user khác
- **Kiểm soát chi tiết**: Có thể cấp quyền theo từng scope (read, write, delete)
- **Bảo mật cao**: Mỗi tài nguyên được bảo vệ riêng biệt

### 2.3. Ví dụ thực tế
- User có thể chia sẻ tài liệu với đồng nghiệp
- Có thể cấp quyền đọc/ghi cho từng tài liệu
- Có thể thu hồi quyền truy cập bất cứ lúc nào

## 3. Requesting Party Token (RPT)

### 3.1. Cấu hình cơ bản
```json
{
  "name": "rpt-policy",
  "type": "rpt",
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

### 3.2. Tác dụng đối với User
- **Token-based Authorization**: Sử dụng token để xác thực quyền truy cập
- **Tính di động**: Token có thể được sử dụng ở nhiều nơi khác nhau
- **Tự động gia hạn**: Token có thể được gia hạn tự động
- **Bảo mật cao**: Token được mã hóa và ký số

### 3.3. Ví dụ thực tế
- User có thể sử dụng token để truy cập API
- Token chứa thông tin về quyền truy cập
- Có thể kiểm tra quyền truy cập mà không cần gọi server

## 4. Kết hợp các phương thức Authorization

### 4.1. Cấu hình kết hợp
```json
{
  "name": "combined-policy",
  "type": "aggregate",
  "policies": [
    {
      "name": "role-policy",
      "type": "role",
      "roles": ["admin"]
    },
    {
      "name": "uma-policy",
      "type": "uma",
      "resources": ["document"]
    },
    {
      "name": "rpt-policy",
      "type": "rpt",
      "permissions": ["read", "write"]
    }
  ],
  "decisionStrategy": "AFFIRMATIVE"
}
```

### 4.2. Tác dụng đối với User
- **Bảo mật nhiều lớp**: Kết hợp nhiều phương thức xác thực
- **Linh hoạt**: Có thể kết hợp các phương thức khác nhau
- **Kiểm soát chi tiết**: Mỗi phương thức có thể kiểm soát một khía cạnh khác nhau
- **Dễ mở rộng**: Có thể thêm các phương thức mới

### 4.3. Ví dụ thực tế
- User phải có role phù hợp
- User phải có quyền truy cập tài nguyên
- User phải có token hợp lệ

## 5. Best Practices

### 5.1. Cấu hình bảo mật
- Sử dụng HTTPS cho tất cả giao tiếp
- Mã hóa token
- Thiết lập thời gian hết hạn token phù hợp
- Sử dụng mật khẩu mạnh

### 5.2. Quản lý quyền
- Phân quyền theo nguyên tắc tối thiểu
- Thường xuyên kiểm tra và cập nhật quyền
- Ghi log đầy đủ các hoạt động
- Có cơ chế khôi phục quyền

### 5.3. Monitoring
- Theo dõi các hoạt động bất thường
- Kiểm tra log thường xuyên
- Có cơ chế cảnh báo
- Backup dữ liệu định kỳ 