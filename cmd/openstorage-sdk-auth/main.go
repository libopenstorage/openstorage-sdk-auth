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
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	yaml "gopkg.in/yaml.v2"

	"github.com/libopenstorage/openstorage-sdk-auth/pkg/auth"
)

var (
	secret      = flag.String("shared-secret", "", "Shared secret to sign token")
	rsaPem      = flag.String("rsa-private-keyfile", "", "RSA Private file to sign token")
	ecdsaPem    = flag.String("ecdsa-private-keyfile", "", "ECDSA Private file to sign token")
	duration    = flag.String("token-duration", "1d", "Duration of time where the token will be valid")
	config      = flag.String("auth-config", "", "Auth configuaration file")
	showVersion = flag.Bool("version", false, "Show version")
	version     = "(dev)"
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(os.Args[0], version)
		return
	}
	if len(*config) == 0 {
		fmt.Println("Must provide a file name")
		os.Exit(1)
	}

	// This is temporary. This program will also support RSA certs
	if len(*secret) == 0 {
		fmt.Println("Must provide a shared secret")
	}
	claims := &auth.Claims{}
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Printf("Failed to read %s: %v", *config, err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, claims); err != nil {
		fmt.Printf("Failed to parse %s: %v", *config, err)
		os.Exit(1)
	}

	token, err := auth.Token(claims,
		&auth.Signature{
			Type: jwt.SigningMethodHS256,
			Key:  []byte(*secret),
		},
		&auth.Options{
			// Temporary
			Expiration: time.Now().Add(time.Minute * 10).Unix(),
		})
	if err != nil {
		fmt.Printf("Failed to create token: %v", err)
		os.Exit(1)
	}
	fmt.Println(token)
}
