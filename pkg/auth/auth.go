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

// Rule provides a method to provide a custom authorization
type Rule struct {
	Allow []string `json:"allow,omitempty" yaml:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty" yaml:"deny,omitempty"`
}

// Claims provides information about the claims in the token
type Claims struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
	Role  string `json:"role" yaml:"role"`
	Rules []Rule `json:"rules" yaml:"rules"`
}

// Signature describes the signature type using definitions from
// the jwt package
type Signature struct {
	Type jwt.SigningMethod
	Key  interface{}
}

// Options provide any options to apply to the token
type Options struct {
	Expiration int64
}

// Token returns a signed JWT containing the claims provided
func Token(
	claims *Claims,
	signature *Signature,
	options *Options,
) (string, error) {

	mapclaims := jwt.MapClaims{
		"email": claims.Email,
		"name":  claims.Name,
		"role":  claims.Role,
		"iat":   time.Now().Unix(),
		"exp":   options.Expiration,
	}
	if claims.Rules != nil {
		mapclaims["rules"] = claims.Rules
	}
	token := jwt.NewWithClaims(signature.Type, mapclaims)
	signedtoken, err := token.SignedString(signature.Key)
	if err != nil {
		return "", err
	}

	return signedtoken, nil
}
