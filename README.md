[![Build Status](https://travis-ci.org/libopenstorage/openstorage-sdk-auth.svg?branch=master)](https://travis-ci.org/libopenstorage/openstorage-sdk-auth)
# OpenStorage SDK Auth

This repo houses the libraries and CLI to create Auth tokens for OpenStorage SDK.

For more information, please see [OpenStorage SDK](https://libopenstorage.github.io)

## Overview
This repo provides the command line program `openstorage-sdk-auth` and Golang package
libraries for users and developers to create auth tokens for OpenStorage SDK.

## Installation

A container will be available, but in the meantime you can do the following:

```
go install -i github.com/libopenstorage/openstorage-sdk-auth/cmd/openstorage-sdk-auth
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
sub: id/luis@portworx.com
email: luis@portworx.com
roles: ["system.user"]
groups: ["px-engineering", "kubernetes-csi"]
```

The yaml has the following structure:
* _email_ string: Email of the account accessing the SDK
* _sub_ string: Unique id of user. Could be the email or a UUID. If this is
  missing, the program will create an ID for the user based on the name and
  email.
* _name_ string: Name of the account accessing the SDK
* _roles_ string list: Roles of the account. This role must already be defined by the
OpenStorage SDK server. The server has the following default roles:
    * system.admin: Access to all APIs
    * system.view: Access to read only APIs only
    * system.user: Access to volume lifecycle APIs only
* _groups_ string list: Groups which the user is part of. Setting the value of `"*"` for the
  group will enable the user of the token to access ALL resources.

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

### Custom Roles
The OpenStorage SDK server allows custom roles. Please see
[OpenStorageRole](https://libopenstorage.github.io/w/master.generated-api.html#serviceopenstorageapiopenstoragerole)
for more information. Once you create a role, you can add it to the token under `roles`.
