package utils

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"fmt"
	"strings"

	"github.com/gruntwork-io/terratest/modules/shell"
)

type IPCounts map[string]int

type Tools struct {
	i *Instance
}

// DigCoreDNS digs CoreDNS for cluster instance
func (t *Tools) DigCoreDNS() []string {
	port := fmt.Sprintf("-p%d", t.i.w.port)
	dnsServer := fmt.Sprintf("@%s", "localhost")
	digApp := shell.Command{
		Command: "dig",
		Args:    []string{port, dnsServer, t.i.GetInfo().Host, "+short", "+tcp", "-4"},
	}
	digAppOut := shell.RunCommandAndGetOutput(t.i.w.t, digApp)
	return strings.Split(digAppOut, "\n")
}

// DigNCoreDNS digs CoreDNS for cluster instance
func (t *Tools) DigNCoreDNS(n int) IPCounts {
	m := make(IPCounts, 0)
	for i := 0; i < n; i++ {
		ips := t.DigCoreDNS()
		for _, ip := range ips {
			m[ip]++
		}
	}
	return m
}

func (t *Tools) Curl() bool {
	return true
}

// IPsHasSimilarProbabilityOnPrecision : ip addresses will appear in the map with a certain probability that is away
// from the average by the deviationPercentage value. For example, we have 400 requests and 4 IP addresses.
// If deviationPercentage =5%, then one address may have 103, the second 97, the third 105 and the fourth 95. Returns true.
// If deviationPercentage =5%, then one address may have 106 hits, the function returns false.
func (f IPCounts) IPsHasSimilarProbabilityOnPrecision(deviationPercentage int) bool {
	var r float64
	for _, v := range f {
		r += float64(v)
	}
	r = r / float64(len(f))
	da := r * float64(100-deviationPercentage) / 100
	db := r * float64(100+deviationPercentage) / 100
	for _, v := range f {
		if float64(v) < da {
			return false
		}
		if float64(v) > db {
			return false
		}
	}
	return true

}
