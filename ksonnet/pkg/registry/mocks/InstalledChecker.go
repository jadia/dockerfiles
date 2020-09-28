// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import pkg "github.com/ksonnet/ksonnet/pkg/pkg"

// InstalledChecker is an autogenerated mock type for the InstalledChecker type
type InstalledChecker struct {
	mock.Mock
}

// IsInstalled provides a mock function with given fields: d
func (_m *InstalledChecker) IsInstalled(d pkg.Descriptor) (bool, error) {
	ret := _m.Called(d)

	var r0 bool
	if rf, ok := ret.Get(0).(func(pkg.Descriptor) bool); ok {
		r0 = rf(d)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(pkg.Descriptor) error); ok {
		r1 = rf(d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
