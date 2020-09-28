// Copyright 2017 The kubecfg authors
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

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/ksonnet/ksonnet/pkg/actions"
	"github.com/ksonnet/ksonnet/pkg/clicmd"
)

// Version is overridden using `-X main.version` during release builds
var version = "(dev build)"
var apimachineryVersion = ""

func main() {
	clicmd.Version = version
	clicmd.APImachineryVersion = apimachineryVersion

	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Error("unable to find working directory")
		os.Exit(1)
	}

	rootCmd, err := clicmd.NewRoot(afero.NewOsFs(), wd, os.Args[1:])
	if err != nil {
		log.WithError(err).Error("unable create ksonnet command")
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		// PersistentPreRunE may not have been run for early
		// errors, like invalid command line flags.
		logFmt := &log.TextFormatter{
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			QuoteEmptyFields:       true,
		}
		log.SetFormatter(logFmt)

		switch err {
		case actions.ErrDiffFound:
			os.Exit(10)
		default:
			log.Error(err.Error())
			os.Exit(1)
		}
	}
}