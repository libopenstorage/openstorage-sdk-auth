/*
Copyright 2018 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

// Claims provides information about the claims in the token
// See https://openid.net/specs/openid-connect-core-1_0.html#IDToken
// for more information.
type Claims struct {
	// Issuer is the token issuer. For selfsigned token do not prefix
	// with `https://`.
	Issuer string `json:"iss"`
	// Subject identifier. Unique ID of this account
	Subject string `json:"sub" yaml:"sub"`
	// Account name
	Name string `json:"name" yaml:"name"`
	// Account email
	Email string `json:"email" yaml:"email"`
	// Roles of this account
	Roles []string `json:"roles,omitempty" yaml:"roles,omitempty"`
	// (optional) Groups in which this account is part of
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// Options provide any options to apply to the token
type Options struct {
	// Expiration time in Unix format as per JWT standard
	Expiration int64
	// Issuer of the claims
	Issuer string
}

// Token returns a signed JWT containing the claims provided
func Token(
	claims *Claims,
	signature *Signature,
	options *Options,
) (string, error) {

	mapclaims := jwt.MapClaims{
		"sub":   claims.Subject,
		"iss":   options.Issuer,
		"email": claims.Email,
		"name":  claims.Name,
		"roles": claims.Roles,
		"iat":   time.Now().Unix(),
		"exp":   options.Expiration,
	}
	if claims.Groups != nil {
		mapclaims["groups"] = claims.Groups
	}
	token := jwt.NewWithClaims(signature.Type, mapclaims)
	signedtoken, err := token.SignedString(signature.Key)
	if err != nil {
		return "", err
	}

	return signedtoken, nil
}
