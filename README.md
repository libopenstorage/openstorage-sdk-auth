[![Build Status](https://travis-ci.org/libopenstorage/openstorage-sdk-auth.svg?branch=master)](https://travis-ci.org/libopenstorage/openstorage-sdk-auth)
# OpenStorage SDK Auth

This repo houses the libraries and CLI to create Auth tokens for OpenStorage SDK.

For more information, please see [OpenStorage SDK](https://libopenstorage.github.io)

## Overview
This repo provides the command line program `openstorage-sdk-auth` and Golang package
libraries for users and developers to create auth tokens for OpenStorage SDK.

OpenStorage Auth tokens are JWTs based loosely based on the OIDC specification.

## Installation

A container will be available, but in the meantime you can do the following:

```
go install -u github.com/libopenstorage/openstorage-sdk-auth/cmd/openstorage-sdk-auth
```

## Usage

To use, you will need to first decide which key type to use to sign the tokens. Although
shared secrets are simple, we recommend using RSA256. In the `tools/` directory you will
find a simple script to generate private and public PEM files.

You will then need to create a claims file using the specification highlighted in this
document. Here is an example of a claims file which defines the email, name, and authorization
of the account:

```yaml
name: Luis Pabon
email: luis@portworx.com
role: volumecreator
groups: ["px-engineering", "kubernetes-csi"]
rules:
  - services: ["volume", "identity"]
    apis: ["*"]
  - services: ["node"]
    apis: ["inspectcurrent"]
```

You can then generate a token using `openstorage-sdk-auth`. In the example below, we generate
a token with an expiration time of 30 days. We use the sample unsecure RSA pem files part
of this repo to sign the token.

```
openstorage-sdk-auth \
  --auth-config=cmd/openstorage-sdk-auth/sample.yml \
  --rsa-private-keyfile=tools/rsa_sample_unsecure_private.pem \
  --token-duration=30d \
  --output=private.token
```

## Token Specification
OpenStorage SDK tokens are JWT tokens whose _claims_ values are highlighted
below:

* _email_: Email of the account accessing the SDK
* _name_: Name of the account accessing the SDK
* _role_: (optional) Role of the account. This role must already be defined by the
OpenStorage SDK server.
* _groups_: (optional) Groups which the user is part of
* _rules_: (optional) Custom role definitions. This allows the token
  to define the APIs and services the account is allowed to access. Please
  see below for more information.
* _exp_: Time when the token will expire
* _iat_: Time when the token was created

### Custom Authorization
The OpenStorage SDK server is equipped to handle customized authorization
claims. Using this model allows the token generator to customize the authorization
rules of the token for specific accounts.

Creating custom authorizations is done by setting up a set of allowed _rules_
directives which are sequentially scanned until a match is found. Rules
are created from the OpenStorage SDK API names as follows:

* services: Is the gRPC service name in `OpenStorage<service name>` in lowercase
* apis: Is the API name in the service in lowercase

Rules can also be set to `*` to allow all services or apis.

Here is an example of a set of rules:

* Allow any call:

```yaml
rules:
  - services: ["*"]
    apis: ["*"]
```

* Allow only cluster operations:

```yaml
rules:
  - services: ["cluster"]
    apis: ["*"]
```

* Allow Identity service, able to inspect only the current node, and ability to only create API
  volume and snapshot calls

```yaml
rules:
  - services: ["identity", "volume"]
    apis: ["create", "version", "capabilities"]
  - services: ["node"]
    apis: ["inspectcurrent"]
```

