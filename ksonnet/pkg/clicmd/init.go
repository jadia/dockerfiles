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

package clicmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/ksonnet/ksonnet/pkg/actions"
	"github.com/ksonnet/ksonnet/pkg/client"
	"github.com/spf13/cobra"
)

const (
	initShortDesc = "Initialize a ksonnet application"

	vInitAPISpec               = "init-api-spec"
	vInitSkipDefaultRegistries = "init-skip-default-registries"
	vInitEnvironment           = "init-environment"
)

var (
	initLong = `
The ` + "`init`" + ` command initializes a ksonnet application in a new directory,` + " `app-name`" + `.

This command generates all the project scaffolding required to begin creating and
deploying components to Kubernetes clusters.

ksonnet applications are initialized based on your current cluster configurations,
as defined in your` + " `$KUBECONFIG` " + `environment variable. The *Examples* section
below demonstrates how to customize these configurations.

Creating a ksonnet application results in the following directory tree.

    app-name/
      .ksonnet/      Metadata for ksonnet
      app.yaml       Application specifications (e.g. name, API version)
      components/    Top-level Kubernetes objects defining the application
      environments/  Kubernetes cluster definitions
        default/     Default environment, initialized from the current kubeconfig
          .metadata/ Contains a versioned ksonnet-lib, see [1] for details
      lib/           User-written .libsonnet files
      vendor/        Libraries that define prototypes and their constituent parts

To begin populating your ksonnet application, see the docs for` + " `ks generate` " + `.

[1] ` + "`ksonnet-lib`" + ` is a Jsonnet helper library that wraps Kubernetes-API-compatible
types. A specific version of ` + "`ksonnet-lib`" + ` is automatically provided for each
environment. Users can set flags to generate the library based on a variety of data,
including server configuration and an OpenAPI specification of a specific Kubernetes
build. By default, this is generated using cluster information specified by the
current context, in the file pointed to by` + " `$KUBECONFIG`" + `.

### Related Commands

* ` + "`ks generate` " + `— ` + protoShortDesc["use"] + `

### Syntax
`
	initExample = `# Initialize a ksonnet application, based on cluster information from the
# active kubeconfig file (as specified by the environment variable $KUBECONFIG).
# More specifically, the current context is used.
ks init app-name

# Initialize a ksonnet application, using the context 'dev' from the current
# kubeconfig file ($KUBECONFIG). The default environment is created using the
# server address and default namespace located at the context 'dev'.
ks init app-name --context=dev

# Initialize a ksonnet application, using the context 'dev' and the namespace
# 'dc-west' from the current kubeconfig file ($KUBECONFIG). The default environment
# is created using the server address from the 'dev' context, and the specified
# 'dc-west' namespace.
ks init app-name --context=dev --namespace=dc-west

# Initialize a ksonnet application, using v1.7.1 of the Kubernetes OpenAPI spec
# to generate 'ksonnet-lib'.
ks init app-name --api-spec=version:v1.7.1

# Initialize a ksonnet application, using the OpenAPI spec generated by a
# specific build of Kubernetes to generate 'ksonnet-lib'.
ks init app-name --api-spec=file:swagger.json

# Initialize a ksonnet application, outputting the application directory into
# the specified 'custom-location'.
ks init app-name --dir=custom-location`
)

func newInitCmd(fs afero.Fs, wd string) *cobra.Command {
	clientConfig := client.NewDefaultClientConfig()

	initCmd := &cobra.Command{
		Use:     "init <app-name>",
		Short:   initShortDesc,
		Long:    initLong,
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if len(args) != 1 {
				return fmt.Errorf("'init' takes a single argument that names the application we're initializing")
			}

			appName := args[0]
			appDir := viper.GetString(flagDir)
			if appDir == wd {
				appDir = ""
			}

			appRoot, err := genKsRoot(appName, wd, appDir)
			if err != nil {
				return err
			}

			clientConfig := client.NewDefaultClientConfig()

			server, namespace, err := resolveEnvFlags(flags, clientConfig)
			if err != nil {
				return err
			}

			specFlag := viper.GetString(vInitAPISpec)
			if specFlag == "" {
				specFlag = clientConfig.GetAPISpec()
			}

			m := map[string]interface{}{
				actions.OptionFs:                    fs,
				actions.OptionName:                  appName,
				actions.OptionNewRoot:               appRoot,
				actions.OptionEnvName:               viper.GetString(vInitEnvironment),
				actions.OptionSpecFlag:              specFlag,
				actions.OptionServer:                server,
				actions.OptionNamespace:             namespace,
				actions.OptionSkipDefaultRegistries: viper.GetBool(vInitSkipDefaultRegistries),
				actions.OptionSkipCheckUpgrade:      true,
			}
			addGlobalOptions(m)

			return runAction(actionInit, m)
		},
	}

	clientConfig.BindClientGoFlags(initCmd)

	// TODO: We need to make this default to checking the `kubeconfig` file.
	initCmd.Flags().String(flagAPISpec, "",
		"Manually specified Kubernetes API version. The corresponding OpenAPI spec is used to generate ksonnet's Kubernetes libraries")
	viper.BindPFlag(vInitAPISpec, initCmd.Flag(flagAPISpec))

	initCmd.Flags().Bool(flagSkipDefaultRegistries, false, "Skip configuration of default registries")
	viper.BindPFlag(vInitSkipDefaultRegistries, initCmd.Flag(flagSkipDefaultRegistries))

	initCmd.Flags().String(flagEnv, "", "Name of initial environment to create")
	viper.BindPFlag(vInitEnvironment, initCmd.Flag(flagEnv))

	return initCmd
}

// genKsRoot determines what the filesystem path to the new ksonnet app will be
// based on the app name, execution directory and cli options provided.
func genKsRoot(appName, ksExecDir, newAppDir string) (string, error) {
	// we either need an app name or a directory to put the app in
	if appName == "" && newAppDir == "" {
		return "", errors.New("invalid application name")
	}

	// if the directory we specified is not relative, just use that path
	if filepath.IsAbs(newAppDir) {
		return newAppDir, nil
	}

	// not having an absolute path to the directory requires knowing CWD
	if ksExecDir == "" {
		return "", errors.New("executing in invalid working directory")
	}

	// if we have a relative dir specified, use that
	if newAppDir != "" {
		return filepath.Abs(filepath.Join(ksExecDir, newAppDir))
	}

	// otherwise (common case) use the CWD/appName as path
	return filepath.Join(ksExecDir, appName), nil
}