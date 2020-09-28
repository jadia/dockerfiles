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

import context "context"
import github "github.com/ksonnet/ksonnet/pkg/util/github"
import go_githubgithub "github.com/google/go-github/github"
import mock "github.com/stretchr/testify/mock"

// GitHub is an autogenerated mock type for the GitHub type
type GitHub struct {
	mock.Mock
}

// CommitSHA1 provides a mock function with given fields: ctx, repo, refSpec
func (_m *GitHub) CommitSHA1(ctx context.Context, repo github.Repo, refSpec string) (string, error) {
	ret := _m.Called(ctx, repo, refSpec)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, github.Repo, string) string); ok {
		r0 = rf(ctx, repo, refSpec)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, github.Repo, string) error); ok {
		r1 = rf(ctx, repo, refSpec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Contents provides a mock function with given fields: ctx, repo, path, sha1
func (_m *GitHub) Contents(ctx context.Context, repo github.Repo, path string, sha1 string) (*go_githubgithub.RepositoryContent, []*go_githubgithub.RepositoryContent, error) {
	ret := _m.Called(ctx, repo, path, sha1)

	var r0 *go_githubgithub.RepositoryContent
	if rf, ok := ret.Get(0).(func(context.Context, github.Repo, string, string) *go_githubgithub.RepositoryContent); ok {
		r0 = rf(ctx, repo, path, sha1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*go_githubgithub.RepositoryContent)
		}
	}

	var r1 []*go_githubgithub.RepositoryContent
	if rf, ok := ret.Get(1).(func(context.Context, github.Repo, string, string) []*go_githubgithub.RepositoryContent); ok {
		r1 = rf(ctx, repo, path, sha1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*go_githubgithub.RepositoryContent)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, github.Repo, string, string) error); ok {
		r2 = rf(ctx, repo, path, sha1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ValidateURL provides a mock function with given fields: u
func (_m *GitHub) ValidateURL(u string) error {
	ret := _m.Called(u)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
