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
	secret   = flag.String("shared-secret", "", "Shared secret to sign token")
	rsaPem   = flag.String("rsa-private-keyfile", "", "RSA Private file to sign token")
	ecdsaPem = flag.String("ecdsa-private-keyfile", "", "ECDSA Private file to sign token")
	duration = flag.String("token-duration", "1d", "Duration of time where the token will be valid. "+
		"Postfix the duration by using "+
		secondDef+" for seconds, "+
		minuteDef+" for minutes, "+
		hourDef+" for hours, "+
		dayDef+" for days, and "+
		yearDef+" for years.")
	config = flag.String("auth-config", "", "Auth account information file providing "+
		"email, name, etc.")
	showVersion = flag.Bool("version", false, "Show version")
	output      = flag.String("output", "", "Output token to file instead of standard out")
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

	// Get claims
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

	// Get duration
	options := &auth.Options{}
	expDuration, err := parseToDuration(*duration)
	if err != nil {
		fmt.Printf("Unable to parse duration")
		os.Exit(1)
	}
	options.Expiration = time.Now().Add(expDuration).Unix()

	// Get signature
	signature := &auth.Signature{}
	if len(*secret) != 0 {
		signature.Key = []byte(*secret)
		signature.Type = jwt.SigningMethodHS256
	} else if len(*rsaPem) != 0 {
		pem, err := ioutil.ReadFile(*rsaPem)
		if err != nil {
			fmt.Printf("Failed to read RSA file: %v", err)
			os.Exit(1)
		}
		signature.Key, err = jwt.ParseRSAPrivateKeyFromPEM(pem)
		if err != nil {
			fmt.Printf("Failed to parse RSA file: %v", err)
			os.Exit(1)
		}
		signature.Type = jwt.SigningMethodRS256
	} else if len(*ecdsaPem) != 0 {
		pem, err := ioutil.ReadFile(*ecdsaPem)
		if err != nil {
			fmt.Printf("Failed to read ECDSA file: %v", err)
			os.Exit(1)
		}
		signature.Key, err = jwt.ParseECPrivateKeyFromPEM(pem)
		if err != nil {
			fmt.Printf("Failed to parse ECDSA file: %v", err)
			os.Exit(1)
		}
		signature.Type = jwt.SigningMethodES256
	} else {
		fmt.Printf("Must provide a secret key to sign token")
		os.Exit(1)
	}

	// Generate token
	token, err := auth.Token(claims, signature, options)
	if err != nil {
		fmt.Printf("Failed to create token: %v", err)
		os.Exit(1)
	}

	// Print token
	if len(*output) != 0 {
		err := ioutil.WriteFile(*output, []byte(token), 0600)
		if err != nil {
			fmt.Printf("Failed to create %s: %v", *output, err)
			os.Exit(1)
		}
	} else {
		fmt.Println(token)
	}
}
