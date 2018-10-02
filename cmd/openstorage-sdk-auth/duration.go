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
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	secondDef = "s"
	minuteDef = "m"
	hourDef   = "h"
	dayDef    = "d"
	yearDef   = "y"

	day  = time.Hour * 24
	year = day * 365
)

var (
	secondRegex = regexp.MustCompile("([0-9]+)" + secondDef)
	minuteRegex = regexp.MustCompile("([0-9]+)" + minuteDef)
	hourRegex   = regexp.MustCompile("([0-9]+)" + hourDef)
	dayRegex    = regexp.MustCompile("([0-9]+)" + dayDef)
	yearRegex   = regexp.MustCompile("([0-9]+)" + yearDef)
)

func parseToDuration(s string) (time.Duration, error) {

	regexs := []struct {
		regex    *regexp.Regexp
		duration time.Duration
	}{
		{
			regex:    secondRegex,
			duration: time.Second,
		},
		{
			regex:    minuteRegex,
			duration: time.Minute,
		},
		{
			regex:    hourRegex,
			duration: time.Hour,
		},
		{
			regex:    dayRegex,
			duration: day,
		},
		{
			regex:    yearRegex,
			duration: year,
		},
	}
	for _, r := range regexs {
		if val := r.regex.FindString(s); len(val) != 0 {
			parsed, err := strconv.Atoi(val[:len(val)-1])
			if err != nil {
				return 0, err
			}
			return time.Duration(parsed) * r.duration, nil
		}
	}

	return 0, fmt.Errorf("Unable to parse")
}
