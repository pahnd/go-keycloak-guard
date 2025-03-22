<!-- TOC -->
  * [Giới thiệu](#about)
  * [Người đóng góp](#contributors)
  * [Yêu cầu](#requirements)
  * [Cấu trúc thư mục](#file-structure)
  * [Quy trình Phương thức Phân quyền](#authorization-methods-workflow)
    * [EnableUMAAuthorization](#enableumaauthorization)
    * [EnableRPTAuthorization](#enablerptauthorization)
    * [EnableRoleBasedAuthorization](#enablerolebasedauthorization)
    * [Quy trình Phân quyền Kết hợp](#combined-authorization-workflow)
    * [Tính năng Chính](#key-features)
    * [Ví dụ Quy trình](#example-workflow)
  * [Cài đặt](#installation)
    * [Biên dịch](#compiling)
    * [Biến Môi trường cần thiết cho cài đặt plugin Kong](#env-variables-required-for-kong-plugin-installation)
    * [schema.lua](#schemalua)
      * [Mô tả](#description)
        * [Ví dụ Vai trò trong Kong:](#example-role-in-kong)
        * [Ví dụ Vai trò trong Konga:](#example-role-in-konga)
        * [Tóm tắt](#summary)
      * [Cài đặt](#installation-1)
    * [Ví dụ cài đặt bổ sung](#additional-installation-examples)
    * [Docker - Kong, Konga & plugin](#docker---kong-konga--the-plugin)
      * [Ví dụ cách cấu hình Kong để sử dụng plugin bằng http requests](#examples-on-how-to-configure-kong-to-use-the-plugin-using-http-requests)
        * [Tạo service](#create-a-service)
        * [Tạo route cho service](#create-a-route-for-that-service)
        * [Kích hoạt plugin keycloak-guard cho Service](#activate-the-keycloak-guard-plugin-for-the-service)
        * [Kích hoạt plugin keycloak-guard cho Route cụ thể](#activate-the-keycloak-guard-plugin-to-a-specific-route)
      * [Ví dụ cách cấu hình Kong để sử dụng plugin qua Konga](#examples-on-how-to-configure-kong-to-use-the-plugin-via-konga)
  * [Cấu hình Keycloak](#keycloak-configuration)
    * [EnableUMAAuthorization](#enableumaauthorization-1)
      * [Ví dụ Tạo User](#create-user-example)
      * [Ví dụ Tạo Scope](#create-scope-example)
      * [Ví dụ Tạo Resource](#create-resource-example)
      * [Ví dụ Tạo Client Policy](#create-client-policy-example)
      * [Tạo Permission](#create-permission)
      * [Tạo Audience scope mapper cho client2](#create-an-audience-scope-mapper-for-client2)
      * [Cách sử dụng](#how-to-use)
        * [Lấy access token cho user:](#fetch-the-access-token-for-the-user)
        * [Sử dụng access token để gửi request qua kong.](#use-the-access-token-to-make-a-request-through-kong-)
    * [EnableRPTAuthorization](#enablerptauthorization-1)
      * [Ví dụ Tạo Scope](#create-scope-example-1)
      * [Ví dụ Tạo Resource](#create-resource-example-1)
      * [Ví dụ Tạo Client Policy](#create-client-policy-example-1)
      * [Tạo Permission](#create-permission-1)
      * [Cách sử dụng](#how-to-use-1)
        * [Cấu hình plugin Kong cho quy trình RPT](#kong-plugin-configuration-for-the-rpt-workflow-)
        * [Yêu cầu permission ticket](#request-the-permission-ticket)
        * [Lấy Client Credentials Token cho "client2"](#obtain-the-client-credentials-token-for-client2)
        * [Sử dụng Client Credentials Token và Permission ticket để lấy RPT Token](#use-the-client-credentials-token-and-the-permission-ticket-to-obtain-the-rpt-token)
        * [Sử dụng RPT token để gọi các API được phân quyền đến client1 (qua Kong)](#use-the-rpt-token-to-make-authorized-calls-to-client1-through-kong)
    * [EnableRoleBasedAuthorization](#enablerolebasedauthorization-1)
      * [Ví dụ Tạo User](#create-user-example-1)
      * [Tạo Audience scope mapper cho client2](#create-an-audience-scope-mapper-for-client2-1)
      * [Tạo vai trò mới cho client1](#create-a-new-role-for-client1)
      * [Gán vai trò mới cho user](#assign-the-newly-created-role-to-a-user)
      * [Cách sử dụng](#how-to-use-2)
        * [Cấu hình plugin Kong cho quy trình Phân quyền dựa trên Vai trò](#kong-plugin-configuration-for-the-role-based-authorization-workflow)
        * [Lấy access token của user (cho user1, client2)](#obtain-user-access-token-for-user1-client2)
        * [Sử dụng Access token vừa lấy được để gọi các API được phân quyền đến client1 (qua Kong)](#use-the-access-token-that-you-have-just-obtained-to-make-authorized-calls-to-client1-through-kong)
  * [Các vấn đề đã biết:](#known-issues)
  * [Luồng hoạt động của hệ thống](#luồng-hoạt-động-của-hệ-thống)
    * [1. Luồng xác thực cơ bản](#1-luồng-xác-thực-cơ-bản)
    * [2. Quá trình xác thực với Keycloak](#2-quá-trình-xác-thực-với-keycloak)
    * [3. Quá trình xác thực nâng cao](#3-quá-trình-xác-thực-nâng-cao)
    * [4. Luồng xử lý kết quả](#4-luồng-xử-lý-kết-quả)
    * [5. Cấu hình Keycloak](#5-cấu-hình-keycloak)
    * [6. Bảo mật](#6-bảo-mật)
    * [7. Monitoring và Debug](#7-monitoring-và-debug)
<!-- TOC -->

## Giới thiệu
Plugin Kong cho Keycloak quản lý cả xác thực và phân quyền cho các yêu cầu API.
Tài liệu này giải thích cách thiết lập và sử dụng plugin Keycloak Guard với Kong và Konga.

![Plugin configuration via konga](docs/resources/konga_setup.png)

## Người đóng góp
- Tên: Mihai Florentin Mihaila
- Website: https://github.com/mihaiflorentin88

## Yêu cầu
Mặc dù có thể hoạt động với các phiên bản khác, đây là các phiên bản tôi đã thử nghiệm plugin:
- Golang: 1.22.4
- Kong: 3.4.2
- Konga: 0.14.9
- Keycloak: 25.0.1 [Docker](https://github.com/eabykov/keycloak-compose)

## Cấu trúc thư mục

```
├── cmd/ - Chứa các điểm vào. Có thể truy cập cả các thành phần domain và infrastructure.
├── docs/ - Tài liệu.
├── domain/ - Chứa các thành phần domain với quy tắc nghiêm ngặt không sử dụng các phụ thuộc bên ngoài.
├── infrastructure/ - Chứa logic cho các client bên ngoài như API hoặc giải pháp lưu trữ.
└── port/ - Chứa Ports(Interfaces)/DTOs.
```

## Quy trình Phương thức Phân quyền
### EnableUMAAuthorization

- Phương thức này xác thực quyền UMA cục bộ.
- Yêu cầu một access token.
- Sử dụng Resource(s) và Scope(s) được cung cấp cùng với một Chiến lược được định nghĩa.
- Chiến lược có thể là affirmative, consensus hoặc unanimous và được sử dụng để xác định cách xác thực quyền.
- Nếu được đặt là true thì các trường sau sẽ trở thành bắt buộc
  - EnableAuth: Xác thực header Authorization Bearer
  - Permissions: Danh sách quyền. Bạn có thể cung cấp quyền theo chuẩn này: ResourceName#ScopeName
  - Strategy: bạn có thể chọn một trong 3 tùy chọn Chiến lược. Tùy chọn này xác định cách xác thực quyền.

### EnableRPTAuthorization
- Phương thức này sử dụng quy trình permission ticket UMA.
- Nếu header Authorization thiếu trong yêu cầu, plugin sẽ trả lời với một permission ticket.
- Người yêu cầu phải chuyển đổi permission ticket này thành RPT token, sau đó được sử dụng để truy cập tài nguyên.
- nếu được đặt là true thì các tùy chọn sau sẽ trở thành bắt buộc
    - EnableAuth: Xác thực header Authorization Bearer
    - ResourceIDs: Danh sách ID tài nguyên. 

### EnableRoleBasedAuthorization

- Phương thức này xác thực xem user có vai trò được chỉ định hay không.
- Yêu cầu trường `Role` phải được chỉ định.
- Không thể được kích hoạt đồng thời với `EnableRPTAuthorization` hoặc `EnableUMAAuthorization`.
- Nếu được đặt là true, các trường sau sẽ trở thành bắt buộc:
    - EnableAuth: Xác thực header Authorization Bearer
    - Role: Vai trò bắt buộc mà user phải có


### Quy trình Phân quyền Kết hợp

Khi cả hai phương thức phân quyền được kích hoạt, plugin ưu tiên quy trình RPT. Đây là cách nó hoạt động:

1. Token Phân quyền Thiếu hoặc Không hợp lệ:
   - Plugin trả lời với một permission ticket.
   - Client có thể trao đổi permission ticket này để lấy RPT token bằng cách sử dụng Keycloak Authorization API hoặc sử dụng access token có quyền phù hợp.
2. Header Phân quyền Hiện diện:
   - Plugin kiểm tra tính hợp lệ của access token trong header Authorization.
   - Nếu access token hợp lệ, nó xác thực xem token có phải là RPT hay không.
   - Nếu token không phải là RPT, plugin chuyển sang quy trình Phân quyền UMA để xác thực quyền dựa trên Resource(s), Scope(s) và Chiến lược được định nghĩa trước.

Đây là ví dụ về response body sẽ được trả về nếu header Authorization Bearer thiếu hoặc không hợp lệ:
```json
{
    "message": "The request is missing the Requesting Party Token (RPT). Please obtain an RPT using the provided permission ticket.",
    "code": 401,
    "permissionTicket": "keycloak-unique-generated-permission-ticket"
}
```

Nếu key permissionTicket hiện diện trong response thì người yêu cầu phải tạo RPT. Hoặc nếu EnableUMAAuthorization được bật thì người yêu cầu cũng có thể cung cấp một access token hợp lệ.
RPT sẽ phải được cung cấp dưới dạng header Authorization Bearer.

### Tính năng Chính

- Xử lý Ưu tiên: Ưu tiên quy trình RPT khi cả hai phương thức được kích hoạt, đảm bảo rằng các client không có token hợp lệ nhận được permission ticket để kiểm soát truy cập động.
- Cơ chế Fallback: Sử dụng Phân quyền UMA làm fallback để xác thực các token không phải RPT, đảm bảo quản lý kiểm soát truy cập toàn diện.
- Tích hợp Liền mạch: Tích hợp cả hai phương thức phân quyền một cách liền mạch để cung cấp cơ chế bảo mật linh hoạt và mạnh mẽ.

### Ví dụ Quy trình
1. Yêu cầu không có Token Phân quyền:
   - Client nhận được permission ticket trong response.
   - Client phải chuyển đổi ticket này thành RPT token để truy cập tài nguyên.
2. Yêu cầu với Access Token Hợp lệ:
   - Plugin kiểm tra token để xác thực tính hợp lệ và loại của nó.
   - Nếu token không phải là RPT, quy trình Phân quyền UMA được sử dụng để xác thực quyền.

## Cài đặt
### Biên dịch
```bash
make compile # Chỉ hoạt động nếu bạn đã cài đặt golang 1.22.4 trên hệ thống
make docker-compile # Sử dụng container docker để biên dịch binary
```
Theo mặc định, nó biên dịch cho linux trên kiến trúc amd64.
Nếu bạn muốn biên dịch cho các nền tảng hoặc kiến trúc khác, sử dụng một trong các lệnh dưới đây (yêu cầu cài đặt golang 1.22.4) hoặc bạn có thể sửa đổi Makefile và sử dụng docker để biên dịch
```bash
# Windows x86 64 bit
go mod tidy && GOOS=windows GOARCH=amd64 go build -o bin//keycloak-guard-windows-amd64.exe main.go
# Windows ARM 64 bit
go mod tidy && GOOS=windows GOARCH=arm64 go build -o bin//keycloak-guard-windows-arm64.exe main.go

# MacOS
# MacOS Darwin x86 64 bit
go mod tidy && GOOS=darwin GOARCH=amd64 go build -o bin//keycloak-guard-darwin-amd64 main.go
# MacOS Darwin ARM 64 bit
go mod tidy && GOOS=darwin GOARCH=arm64 go build -o bin//keycloak-guard-darwin-arm64 main.go

# Linux
# Linux x86 64 bit
go mod tidy && GOOS=linux GOARCH=amd64 go build -o bin//keycloak-guard-linux-amd64 main.go
# Linux ARM 64 bit
go mod tidy && GOOS=linux GOARCH=arm64 go build -o bin//keycloak-guard-linux-arm64 main.go
```

### Biến Môi trường cần thiết cho cài đặt plugin Kong
```bash
# Giả định rằng tên binary của bạn là keycloak-guard
export KONG_PLUGINSERVER_NAMES="keycloak-guard"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_START_CMD="/usr/bin/keycloak-guard -kong-prefix /tmp"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_QUERY_CMD="/usr/bin/keycloak-guard -dump"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_SOCKET="/tmp/keycloak-guard.socket"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_START_TIMEOUT="10"
export KONG_PLUGINS="bundled,keycloak-guard"
```

### schema.lua

#### Mô tả

schema.lua là một script Lua được sử dụng trong các plugin Kong để định nghĩa schema cấu hình cho plugin. Nó đóng vai trò quan trọng trong việc xác thực và quản lý cài đặt cấu hình của plugin. Đây là giải thích ngắn gọn về vai trò của nó:
1. Định nghĩa Cấu trúc Cấu hình: schema.lua chỉ định cấu trúc của các tùy chọn cấu hình mà người dùng có thể đặt cho plugin. Điều này bao gồm việc định nghĩa các trường, kiểu dữ liệu của chúng, giá trị mặc định và quy tắc xác thực.
2. Đảm bảo Tính hợp lệ: Nó đảm bảo rằng cấu hình được cung cấp bởi người dùng là hợp lệ và đáp ứng các tiêu chí mong đợi trước khi plugin được thực thi. Việc xác thực này giúp ngăn chặn lỗi runtime do cấu hình không chính xác.
3. Tích hợp với Konga: Khi sử dụng Konga, một giao diện người dùng để quản lý Kong, schema.lua giúp Konga hiểu các tùy chọn cấu hình có sẵn cho plugin, cho phép giao diện thân thiện với người dùng để thiết lập và sửa đổi các tùy chọn này.

##### Ví dụ Vai trò trong Kong:

1. Định nghĩa Trường: Chỉ định các trường như api_key, timeout và kiểu dữ liệu tương ứng của chúng (string, number, v.v.).
2. Xác thực: Thực thi các quy tắc như trường bắt buộc, độ dài trường và phạm vi giá trị chấp nhận được.
3. Giá trị Mặc định: Cung cấp giá trị mặc định cho cài đặt cấu hình nếu người dùng không chỉ định chúng.

##### Ví dụ Vai trò trong Konga:

Tích hợp Giao diện Người dùng: Cho phép Konga tạo động các biểu mẫu và trường nhập liệu dựa trên schema của plugin, cho phép người dùng cấu hình plugin thông qua giao diện Konga một cách dễ dàng.

##### Tóm tắt
Tóm lại, schema.lua là cần thiết để định nghĩa, xác thực và quản lý cấu hình của các plugin Kong, đảm bảo tích hợp và hoạt động mượt mà trong cả môi trường Kong và Konga.

#### Cài đặt
Sao chép file schema.lua từ thư mục gốc của repository đến đường dẫn này: ```/usr/local/share/lua/5.1/kong/plugins/keycloak-guard/schema.lua```

### Ví dụ cài đặt bổ sung
Bạn có thể tìm thêm ví dụ về cách thiết lập plugin và [schema.lua](./schema.lua) trong file [docker-compose.yaml](./docker-compose.yaml).

### Docker - Kong, Konga & plugin
Repository bao gồm file [docker-compose.yaml](./docker-compose.yaml) thiết lập môi trường hoàn chỉnh với Kong, Konga và plugin tùy chỉnh được cài đặt. Để quản lý các service này, bạn có thể sử dụng các lệnh Makefile được cung cấp:
```bash
make kong-start # Khởi động các container Kong và Konga
make kong-stop # Dừng các container Kong và Konga
make docker-clean-up # Dừng các container, xóa tất cả images và networks
```
#### Ví dụ cách cấu hình Kong để sử dụng plugin bằng http requests
Để tạo service, thêm route và gán plugin keycloak-guard trong Kong, bạn có thể sử dụng các lệnh curl sau:
##### Tạo service
```bash
curl -i -X POST http://localhost:8001/services/ \
  --data name=example-service \
  --data url=http://your.service
```
##### Tạo route cho service
```bash
curl -i -X POST http://localhost:8001/services/example-service/routes \
  --data 'paths[]=/example'
```
##### Kích hoạt plugin keycloak-guard cho Service
```bash
curl -i -X POST http://localhost:8001/services/example-service/plugins \
  --data name=keycloak-guard \
  --data config.KeycloakURL=http://your-keycloak-url \
  --data config.Realm=your-realm \
  --data config.ClientID=your-client-id \
  --data config.ClientSecret=your-client-secret \
  --data config.EnableAuth=true \ # Tùy chọn nếu EnableUMAAuthorization và EnableRPTAuthorization được đặt là false 
  --data config.EnableUMAAuthorization=true \ # Tùy chọn
  --data config.Permissions[]=resouceName#exampleScope \ # Tùy chọn nếu EnableUMAAuthorization được đặt là false
  --data config.Strategy=affirmative \ # Tùy chọn nếu EnableUMAAuthorization được đặt là false
  --data config.EnableRPTAuthorization=true \ # Tùy chọn
  --data config.ResourceIDs[]=resource-id-1 \ # Tùy chọn nếu EnableRPTAuthorization được đặt là false
  --data config.ResourceIDs[]=resource-id-2 # # Tùy chọn nếu EnableRPTAuthorization được đặt là false
  --data config.EnableRoleBasedAuthorization=true # Tùy chọn nếu EnableRoleBasedAuthorization được đặt là true thì EnableAuth phải được đặt là true và Role phải được thêm vào. Đồng thời cả EnableUMAAuthorization và EnableRPTAuthorization phải được đặt là false
  --data config.Role=role1 # Tùy chọn yêu cầu EnableRoleBasedAuthorization được đặt là true
```

##### Kích hoạt plugin keycloak-guard cho Route cụ thể
```bash
curl -i -X POST http://localhost:8001/routes/{route_id}/plugins \
  --data name=keycloak-guard \
  --data config.KeycloakURL=http://your-keycloak-url \
  --data config.Realm=your-realm \
  --data config.ClientID=your-client-id \
  --data config.ClientSecret=your-client-secret \
  --data config.EnableAuth=true \ # Tùy chọn nếu EnableUMAAuthorization và EnableRPTAuthorization được đặt là false 
  --data config.EnableUMAAuthorization=true \ # Tùy chọn
  --data config.Permissions[]=resouceName#exampleScope \ # Tùy chọn nếu EnableUMAAuthorization được đặt là false
  --data config.Strategy=affirmative \ # Tùy chọn nếu EnableUMAAuthorization được đặt là false
  --data config.EnableRPTAuthorization=true \ # Tùy chọn
  --data config.ResourceIDs[]=resource-id-1 \ # Tùy chọn nếu EnableRPTAuthorization được đặt là false
  --data config.ResourceIDs[]=resource-id-2 # # Tùy chọn nếu EnableRPTAuthorization được đặt là false
  --data config.EnableRoleBasedAuthorization=true # Tùy chọn nếu EnableRoleBasedAuthorization được đặt là true thì EnableAuth phải được đặt là true và Role phải được thêm vào. Đồng thời cả EnableUMAAuthorization và EnableRPTAuthorization phải được đặt là false
  --data config.Role=role1 # Tùy chọn yêu cầu EnableRoleBasedAuthorization được đặt là true
```
#### Ví dụ cách cấu hình Kong để sử dụng plugin qua Konga
[Ảnh chụp màn hình](docs/resources/konga_setup.png) này chứa ví dụ về cách thiết lập plugin qua Konga với tất cả các tính năng được bật.
![Plugin configuration via konga](docs/resources/konga_setup.png)

## Cấu hình Keycloak
Tôi sẽ sử dụng "client1" và "client2" trong các ví dụ của tôi. "client2" sẽ là client gửi các yêu cầu đến client1
### EnableUMAAuthorization
Để quy trình này hoạt động, bạn sẽ cần có cấu hình Keycloak sau:
1. **client2:**
   - tạo audience scope mapper
2. **Realm:**
   - tạo user đang hoạt động hoặc sử dụng user hiện có.
3. **client1:**
   - một scope mới
   - một resource mới
   - một client policy mới 
   - một permission mới 

#### Ví dụ Tạo User
Tạo user:
![Create user](docs/resources/keycloak_create_user.png)
Đặt mật khẩu user
![Set user password](docs/resources/keycloak_set_user_password.png)
#### Ví dụ Tạo Scope
![Create Scope](docs/resources/keycloak_create_scope.png)
#### Ví dụ Tạo Resource
![Create Resource](docs/resources/keycloak_create_resource.png)
#### Ví dụ Tạo Client Policy
![Create Client Policy](docs/resources/keycloak_create_client_policy.png)
#### Tạo Permission
![Create Permission](docs/resources/keycloak_create_permission.png)
#### Tạo Audience scope mapper cho client2
![Audience Scope Mapper](docs/resources/keycloak_audience_scope_mapper.png)
#### Cách sử dụng
##### Lấy access token cho user:
```bash
curl --location 'http://keycloak-url/realms/<realm>/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=client2' \
--data-urlencode 'username=user1' \
--data-urlencode 'password=userpassword' \
--data-urlencode 'client_secret=YOie4lyoXakCXDuP7jRCsUM4Xx4OxUOB'
```
##### Sử dụng access token để gửi request qua kong. 
Plugin keycloak-gateway phải được kích hoạt và cấu hình cho quy trình này.
![UMA Authorization workflow plugin configuration](docs/resources/kong_plugin_configuration_for_uma_workflow.png)
Để biết thêm chi tiết về thiết lập và hướng dẫn sử dụng, tham khảo phần [Ví dụ cài đặt bổ sung](#additional-installation-examples).
```bash
curl --location 'http://kong-hostname:8000/test' \
--header 'Authorization: Bearer <accessToken>'
```
### EnableRPTAuthorization
Để quy trình này hoạt động, bạn sẽ cần có cấu hình Keycloak sau:
1. **client1:**
   - một scope mới
   - một resource mới
   - một client policy mới
   - một permission mới

#### Ví dụ Tạo Scope
![Create Scope](docs/resources/keycloak_create_scope.png)
#### Ví dụ Tạo Resource
![Create Resource](docs/resources/keycloak_create_resource.png)
#### Ví dụ Tạo Client Policy
![Create Client Policy](docs/resources/keycloak_create_client_policy.png)
#### Tạo Permission
![Create Permission](docs/resources/keycloak_create_permission.png)
#### Cách sử dụng
##### Cấu hình plugin Kong cho quy trình RPT 
![RPT Workflow configuration](docs/resources/kong_plugin_configuration_for_rpt_workflow.png)
##### Yêu cầu permission ticket
Bất kỳ cuộc gọi nào không có RPT token đang hoạt động sẽ trả về response chứa permission ticket

```bash
curl --location 'http://kongHostname:8000/test'
```

Response status code sẽ là 401 và response body sẽ như trong ví dụ dưới đây:
```json
{
    "message": "The request is missing the Requesting Party Token (RPT). Please obtain an RPT using the provided permission ticket.",
    "code": 401,
    "permissionTicket": "<permission_ticket>"
}
```
##### Lấy Client Credentials Token cho "client2"
```bash
curl --location 'http://keycloak-url/realms/<realm>/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=client_credentials' \
--data-urlencode 'client_id=client2' \
--data-urlencode 'client_secret=YOie4lyoXakCXDuP7jRCsUM4Xx4OxUOB'
```
Response sẽ chứa Client Credentials token (access_token):
```json
{
    "access_token": "<client_credentials_token>",
    "expires_in": 86400,
    "refresh_expires_in": 0,
    "token_type": "Bearer",
    "not-before-policy": 0,
    "scope": "email profile"
}
```

##### Sử dụng Client Credentials Token và Permission ticket để lấy RPT Token
```bash
curl --location 'http://keycloak-url/realms/<realm>/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Authorization: Bearer <client_credentials_token>' \
--data-urlencode 'grant_type=urn:ietf:params:oauth:grant-type:uma-ticket' \
--data-urlencode 'ticket=<permission_ticket>'
```
Response sẽ chứa RPT token sẽ được sử dụng để phân quyền các cuộc gọi tiếp theo đến ứng dụng của chúng ta (qua kong).
```json
{
    "upgraded": false,
    "access_token": "<rpt_token>",
    "expires_in": 86400,
    "refresh_expires_in": 0,
    "token_type": "Bearer",
    "not-before-policy": 0
}
```
##### Sử dụng RPT token để gọi các API được phân quyền đến client1 (qua Kong)
```bash
curl --location 'http://kongHostname:8000/test' \
--header 'Authorization: Bearer <rpt_token>'
```

### EnableRoleBasedAuthorization
Để quy trình này hoạt động, bạn sẽ cần có cấu hình Keycloak sau:
1. **Realm:**
   - tạo user đang hoạt động hoặc sử dụng user hiện có.
2. **client1:**
   - tạo vai trò mới
   - gán vai trò mới cho username
3. **client2:**
   - tạo audience scope mapper

#### Ví dụ Tạo User
Tạo user:
![Create user](docs/resources/keycloak_create_user.png)
Đặt mật khẩu user
![Set user password](docs/resources/keycloak_set_user_password.png)

#### Tạo Audience scope mapper cho client2
![Audience Scope Mapper](docs/resources/keycloak_audience_scope_mapper.png)

#### Tạo vai trò mới cho client1
![Create a role](docs/resources/keycloak_create_a_role.png)

#### Gán vai trò mới cho user
Tab ánh xạ vai trò cho user:
![Role Mapping](docs/resources/keycloak_user_role_mapping_tab.png)
Gán vai trò cho user đó:
![Role Mapping](docs/resources/keycloak_assign_role_to_user.png)

#### Cách sử dụng
##### Cấu hình plugin Kong cho quy trình Phân quyền dựa trên Vai trò
- EnableAuth phải được bật
- EnableUMAAuthorization phải được tắt
- EnableRPTAuthorization phải được tắt
- EnableRoleBasedAuthorization phải được bật
- Role (trường) phải chứa tên cho vai trò bạn muốn cho phép truy cập (role1)
![Kong Role Based CFG](docs/resources/kong_plugin_configuration_for_role_based_workflow.png)

##### Lấy access token của user (cho user1, client2)
```bash
curl --location 'http://keycloak-url/realms/<realm>/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=client2' \
--data-urlencode 'username=user1' \
--data-urlencode 'password=userpassword' \
--data-urlencode 'client_secret=YOie4lyoXakCXDuP7jRCsUM4Xx4OxUOB'
```
Response sẽ chứa access token:
```json
{
    "access_token": "<access_token>",
    "expires_in": 86400,
    "refresh_expires_in": 86400,
    "refresh_token": "<refresh_token>",
    "token_type": "Bearer",
    "not-before-policy": 0,
    "session_state": "9f0c033d-d7cf-4b1b-ab64-d591ec04edc2",
    "scope": "email profile"
}
```

##### Sử dụng Access token vừa lấy được để gọi các API được phân quyền đến client1 (qua Kong)

```bash
curl --location 'http://kongHostname:8000/test' \
--header 'Authorization: Bearer <access_token>'
```
Nếu user (user1) có vai trò yêu cầu (role1) thì plugin sẽ cho phép bạn thực hiện yêu cầu.

## Các vấn đề đã biết:
1. Ứng dụng Konga sẽ không hiển thị thông tin "details" cho các trường của plugin (Đây là vấn đề với Konga không được duy trì)

## Luồng hoạt động của hệ thống

### 1. Luồng xác thực cơ bản
- Khi user gửi request đến một endpoint được bảo vệ bởi Kong Gateway
- Kong Gateway sẽ chuyển request đến plugin `keycloak-guard`
- Plugin sẽ kiểm tra token trong header của request
- Nếu không có token hoặc token không hợp lệ, user sẽ được chuyển hướng đến trang đăng nhập của Keycloak

### 2. Quá trình xác thực với Keycloak
- User được chuyển hướng đến Keycloak (port 8080)
- User đăng nhập với tài khoản được cấu hình trong Keycloak
- Sau khi đăng nhập thành công, Keycloak sẽ:
  - Tạo JWT token chứa thông tin user
  - Chuyển hướng user về lại endpoint ban đầu với token
  - Token này sẽ được sử dụng cho các request tiếp theo

### 3. Quá trình xác thực nâng cao
Plugin hỗ trợ 3 loại xác thực nâng cao:

#### a. Role-Based Authorization
- Kiểm tra role của user trong token
- So sánh với các role được cấu hình trong plugin
- Chỉ cho phép truy cập nếu user có role phù hợp

#### b. UMA (User-Managed Access)
- Kiểm tra quyền truy cập tài nguyên
- Sử dụng RPT (Requesting Party Token)
- Kiểm tra các permission được cấu hình

#### c. RPT Authorization
- Kiểm tra RPT token
- Xác thực các permission được cấp
- Đảm bảo token còn hiệu lực

### 4. Luồng xử lý kết quả
- Nếu tất cả các kiểm tra xác thực thành công:
  - Request được chuyển tiếp đến backend service
  - Thông tin user được thêm vào header của request
  - Log được ghi lại cho mục đích audit

- Nếu xác thực thất bại:
  - Trả về lỗi 401 (Unauthorized) hoặc 403 (Forbidden)
  - Log chi tiết lý do từ chối truy cập
  - User có thể được chuyển hướng để đăng nhập lại

### 5. Cấu hình Keycloak
- Keycloak chạy trên port 8080
- Sử dụng PostgreSQL làm database
- Có các tính năng:
  - Health check
  - Metrics
  - Preview features
- Admin credentials:
  - Username: admin
  - Password: admin

### 6. Bảo mật
- Tất cả giao tiếp giữa các service đều trong mạng nội bộ (kong-network)
- Keycloak sử dụng HTTPS
- Token được mã hóa và ký số
- Có cơ chế refresh token
- Log đầy đủ cho mục đích audit

### 7. Monitoring và Debug
- Kong Gateway có log level debug
- Keycloak có health check endpoint
- Có thể theo dõi metrics của Keycloak
- Log được ghi ra stdout/stderr
