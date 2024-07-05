<!-- TOC -->
  * [About](#about)
  * [Author](#author)
  * [Contributors](#contributors)
  * [Requirements](#requirements)
  * [File structure](#file-structure)
  * [Authorization Methods Workflow](#authorization-methods-workflow)
    * [EnableUMAAuthorization](#enableumaauthorization)
    * [EnableRPTAuthorization](#enablerptauthorization)
    * [Combined Authorization Workflow](#combined-authorization-workflow)
    * [Key Features](#key-features)
    * [Example Workflow](#example-workflow)
  * [Installation](#installation)
    * [Compiling](#compiling)
    * [ENV Variables required for Kong plugin installation](#env-variables-required-for-kong-plugin-installation)
    * [schema.lua](#schemalua)
      * [Description](#description)
        * [Example Role in Kong:](#example-role-in-kong)
        * [Example Role in Konga:](#example-role-in-konga)
        * [Summary](#summary)
      * [Installation](#installation-1)
    * [Additional installation examples](#additional-installation-examples)
    * [Docker - Kong, Konga & the plugin](#docker---kong-konga--the-plugin)
      * [Examples on how to configure Kong to use the plugin using http requests](#examples-on-how-to-configure-kong-to-use-the-plugin-using-http-requests)
        * [Create a service](#create-a-service)
        * [Create a route for that service](#create-a-route-for-that-service)
        * [Activate the keycloak-guard Plugin for the Service](#activate-the-keycloak-guard-plugin-for-the-service)
        * [Activate the keycloak-guard Plugin to a Specific Route](#activate-the-keycloak-guard-plugin-to-a-specific-route)
      * [Examples on how to configure Kong to use the plugin via Konga](#examples-on-how-to-configure-kong-to-use-the-plugin-via-konga)
<!-- TOC -->

## About
Kong plugin for Keycloak that manages both authentication and authorization for API requests.
This document explains how to set up and use the Keycloak Guard plugin with Kong and Konga.
For more information on the authorization methods, see the [Authorization Methods Workflow](#authorization-methods-workflow) section.

## Author
- Name: Mihai Florentin Mihaila
- Website: https://github.com/mihaiflorentin88

## Contributors
- Name: Mihai Florentin Mihaila
- Website: https://github.com/mihaiflorentin88

## Requirements
While it might work with other versions these are the versions i have tested the plugin with:
- Golang: 1.22.4
- Kong: 3.4.2
- Konga: 0.14.9

## File structure

```
├── cmd/ - Contains entry points. Can access both domain and infrastructure components.
├── docs/ - Documentation resources.
├── domain/ - Contains domain components with the strict rule of never using external dependecies.
├── infrastructure/ - Contains logic for external clients like APIs or storage solutions.
└── port/ - Contains Ports(Interfaces)/DTOs.
```

## Authorization Methods Workflow
### EnableUMAAuthorization

- This method verifies UMA permissions locally.
- It requires an access token.
- It uses the provided Resource(s) and Scope(s) along with a defined Strategy.
- The Strategy can be affirmative, consensus, or unanimous and is used to determine how permissions are validated.
- if set to true then the following fields will be made mandatory
  - EnableAuth: Verifies the Authorization Bearer header
  - Permissions: List of permissions. You can provide the permissions following this standard: ResourceName#ScopeName
  - Strategy: you can choose from one of these 3 options Strategy. This option determines how permissions are validated.

### EnableRPTAuthorization
- This method utilizes the UMA permission ticket workflow.
- If the Authorization header is missing in the request, the plugin responds with a permission ticket.
- The requester must convert this permission ticket into an RPT token, which is then used to gain access to the resource.
- if set to true then the following options will be made mandatory
    - EnableAuth: Verifies the Authorization Bearer header
    - ResourceIDs: List of resource ids. The resource ids are actually Keycloak Permissions that use Resources and Policies in order to determine if the client is authorized. You can use the name of the Permissions for this option.

### Combined Authorization Workflow

When both authorization methods are enabled, the plugin prioritizes the RPT workflow. Here’s how it operates:

1. Authorization Token Missing or Invalid:
   - The plugin responds with a permission ticket.
   - The client must exchange this permission ticket for an RPT token using the Keycloak Authorization API.
2. Authorization Header Present:
   - The plugin checks the validity of the access token in the Authorization header.
   - If the access token is valid, it verifies whether the token is an RPT.
   - If the token is not an RPT, the plugin falls back to the UMA Authorization workflow to validate permissions based on the predefined Resource(s), Scope(s), and Strategy.

This is an example for the response body that will be returned if the Authorization Bearer header is missing or invalid:
```json
{
    "message": "The request is missing the Requesting Party Token (RPT). Please obtain an RPT using the provided permission ticket.",
    "code": 401,
    "permissionTicket": "keycloak-unique-generated-permission-ticket"
}
```

If the permissionTicket key is present in the response then the requester has to generate an RPT. Or if the EnableUMAAuthorization is turned on then the requester will have to provide a valid access token.
The RPT will have to be provided as an Authorization Bearer header.

### Key Features

- Priority Handling: Prioritizes the RPT workflow when both methods are enabled, ensuring that clients without valid tokens receive a permission ticket for dynamic access control.
- Fallback Mechanism: Uses UMA Authorization as a fallback to validate non-RPT tokens, ensuring comprehensive access control management.
- Seamless Integration: Integrates both authorization methods seamlessly to provide flexible and robust security mechanisms.

### Example Workflow
1. Request Without Authorization Token:
   - The client receives a permission ticket in the response.
   - The client must convert this ticket into an RPT token to access the resource.
2. Request With Valid Access Token:
   - The plugin introspects the token to check its validity and type.
   - If the token is not an RPT, the UMA Authorization workflow is used to validate permissions.

## Installation
### Compiling
```bash
make compile # This only works if you have golang 1.22.4 installed on your system
make docker-compile # This uses a docker container to compile the binary
```
By default, it compiles for linux on amd64 architecture.
If you wish to compile for other platforms or architectures use one of the commands bellow (requires golang 1.22.4 installed) or you can modify the Makefile and use docker to compile it
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

### ENV Variables required for Kong plugin installation
```bash
# This assumes that the name of your binary is keycloak-guard
export KONG_PLUGINSERVER_NAMES="keycloak-guard"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_START_CMD="/usr/bin/keycloak-guard -kong-prefix /tmp"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_QUERY_CMD="/usr/bin/keycloak-guard -dump"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_SOCKET="/tmp/keycloak-guard.socket"
export KONG_PLUGINSERVER_KEYCLOAK_GUARD_START_TIMEOUT="10"
export KONG_PLUGINS="bundled,keycloak-guard"
```

### schema.lua

#### Description

schema.lua is a Lua script used in Kong plugins to define the configuration schema for the plugin. It plays a crucial role in validating and managing the plugin’s configuration settings. Here’s a concise explanation of its role:
1. Define Configuration Structure: schema.lua specifies the structure of the configuration options that users can set for the plugin. This includes defining fields, their types, default values, and validation rules.
2. Ensure Validity: It ensures that the configuration provided by the user is valid and meets the expected criteria before the plugin is executed. This validation helps prevent runtime errors due to incorrect configurations.
3. Integration with Konga: When using Konga, a UI for managing Kong, schema.lua helps Konga understand the configuration options available for the plugin, allowing for a user-friendly interface to set and modify these options.

##### Example Role in Kong:

1. Field Definitions: Specifies fields such as api_key, timeout, and their respective data types (string, number, etc.).
2. Validation: Enforces rules like required fields, field length, and acceptable value ranges.
3. Defaults: Provides default values for configuration settings if the user does not specify them.

##### Example Role in Konga:

UI Integration: Enables Konga to dynamically generate forms and input fields based on the plugin’s schema, allowing users to configure the plugin through the Konga interface easily.

##### Summary
In summary, schema.lua is essential for defining, validating, and managing the configuration of Kong plugins, ensuring smooth integration and functionality within both Kong and Konga environments.

#### Installation
Copy the schema.lua file from the repository root to this path: ```/usr/local/share/lua/5.1/kong/plugins/keycloak-guard/schema.lua```

### Additional installation examples
You can find more examples on how to setup the plugin and the [schema.lua](./schema.lua) inside the [docker-compose.yaml](./docker-compose.yaml) file.

### Docker - Kong, Konga & the plugin
The repository includes a [docker-compose.yaml](./docker-compose.yaml) file that sets up a fully functional environment with Kong, Konga, and the custom plugin installed. To manage these services, you can use the provided Makefile commands:
```bash
make kong-start # Start the Kong and Konga containers
make kong-stop # Stop the Kong and Konga containers
make docker-clean-up # Stop the containers, remove all images and networks
```
#### Examples on how to configure Kong to use the plugin using http requests
To create a service, add a route, and assign the keycloak-guard plugin in Kong, you can use the following curl commands:
##### Create a service
```bash
curl -i -X POST http://localhost:8001/services/ \
  --data name=example-service \
  --data url=http://your.service
```
##### Create a route for that service
```bash
curl -i -X POST http://localhost:8001/services/example-service/routes \
  --data 'paths[]=/example'
```
##### Activate the keycloak-guard Plugin for the Service
```bash
curl -i -X POST http://localhost:8001/services/example-service/plugins \
  --data name=keycloak-guard \
  --data config.KeycloakURL=http://your-keycloak-url \
  --data config.Realm=your-realm \
  --data config.ClientID=your-client-id \
  --data config.ClientSecret=your-client-secret \
  --data config.EnableAuth=true \ # Optional if EnableUMAAuthorization and EnableRPTAuthorization are set to false 
  --data config.EnableUMAAuthorization=true \ # Optional
  --data config.Permissions[]=resouceName#exampleScope \ # Optional if EnableUMAAuthorization is set to false
  --data config.Strategy=affirmative \ # Optional if EnableUMAAuthorization is set to false
  --data config.EnableRPTAuthorization=true \ # Optional
  --data config.ResourceIDs[]=resource-id-1 \ # Optional if EnableRPTAuthorization is set to false
  --data config.ResourceIDs[]=resource-id-2 # # Optional if EnableRPTAuthorization is set to false
```

##### Activate the keycloak-guard Plugin to a Specific Route
```bash
curl -i -X POST http://localhost:8001/routes/{route_id}/plugins \
  --data name=keycloak-guard \
  --data config.KeycloakURL=http://your-keycloak-url \
  --data config.Realm=your-realm \
  --data config.ClientID=your-client-id \
  --data config.ClientSecret=your-client-secret \
  --data config.EnableAuth=true \ # Optional if EnableUMAAuthorization and EnableRPTAuthorization are set to false 
  --data config.EnableUMAAuthorization=true \ # Optional
  --data config.Permissions[]=resouceName#exampleScope \ # Optional if EnableUMAAuthorization is set to false
  --data config.Strategy=affirmative \ # Optional if EnableUMAAuthorization is set to false
  --data config.EnableRPTAuthorization=true \ # Optional
  --data config.ResourceIDs[]=resource-id-1 \ # Optional if EnableRPTAuthorization is set to false
  --data config.ResourceIDs[]=resource-id-2 # # Optional if EnableRPTAuthorization is set to false
```
#### Examples on how to configure Kong to use the plugin via Konga
The screenshot contains an example on how to setup the plugin via Konga with all features toggled on.
![](/Users/mihai.mihaila/Workspace/go-kong/docs/resources/konga_setup.png)
