# OpenStorage SDK Auth

This repo houses the libraries and CLI to create Auth tokens for OpenStorage SDK.

For more information, please see [OpenStorage SDK](https://libopenstorage.github.io)

## Overview
This repo provides the command line program `openstorage-sdk-auth` and Golang package
libraries for users and developers to create auth tokens for OpenStorage SDK.

OpenStorage Auth tokens are JWTs based loosely based on the OIDC specification.

## Token Specification
OpenStorage SDK tokens are JWT tokens whose _claims_ values are highlighted
below:

* _email_: Email of the account accessing the SDK
* _name_: Name of the account accessing the SDK
* _role_: Role of the account. This role must already be defined by the
OpenStorage SDK server.
* _exp_: Time when the token will expire
* _iat_: Time when the token was created
* _iss_: TBD
* _sub_: TBD
* _rules_: (optional) Custom role definitions. This allows the token
  to define the APIs and services the account is allowed to access. Please
  see TBD for more information.

### Custom Authorization
The OpenStorage SDK server is equipped to handle customized authorization
claims. Using this model allows the token generator to customize the authorization
rules of the token for specific accounts.

Creating custom authorizations is done by setting up a set of _rules_ using `allow`
or `deny` directives which are sequentially scanned until a match is found. Rules
are applied to the following OpenStorage SDK objects:

* Services: `openstorage.api.<service name>`
* APIs: `openstorage.api.<service name>/<API Name>`

Rules can also be set to `all` to allow or deny any other value.

Here is an example of a set of rules:

* Allow any call:

```json
{
	"rules" : [
		{
			"allow": [
				"all"
			]
		}
	]
}
```

* Allow only cluster operations:

```json
{
	"rules" : [
		{
			"allow": [
				"openstorage.api.OpenStorageCluster"
			]
		},
		{
			"deny": [
				"all"
			]
		}
	]
}
```

* Allow only all API volume calls, with the exception of being unable delete any:

```json
{
	"rules" : [
		{
			"deny": [
				"openstorage.api.OpenStorageVolume/Delete"
			]
		}
		{
			"allow": [
				"openstorage.api.OpenStorageVolume"
			]
		},
		{
			"deny": [
				"all"
			]
		}
	]
}
```

