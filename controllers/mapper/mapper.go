package mapper

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

type Result int

const RecordTypeA = "A"

const (
	ResultExists Result = 1 << iota
	ResultNotFound
	ResultError
	ResultExistsButNotAnnotationFound
	ResultFinalizerRemoved
	ResultFinalizerInstalled
	ResultContinue
)

func (r Result) IsIn(m ...Result) bool {
	for _, v := range m {
		if v == r {
			return true
		}
	}
	return false
}

// Mapper is wrapper around resource. Mappers are an only way to access resources
type Mapper interface {
	Equal(*LoopState) bool
	GetStatus() Status
	GetExposedIPs() ([]string, error)
	TryInjectFinalizer() (Result, error)
	TryRemoveFinalizer(func(*LoopState) error) (Result, error)
	SetReference(*LoopState)
	UpdateStatusAnnotation() error
	RemoveDNSEndpoint() (Result, error)
}
