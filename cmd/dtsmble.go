// Copyright 2023 Interlynk.io
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/interlynk-io/sbomasm/pkg/assemble"
	"github.com/interlynk-io/sbomasm/pkg/dtassemble"
	"github.com/interlynk-io/sbomasm/pkg/logger"
	"github.com/spf13/cobra"
)

// assembleCmd represents the assemble command
var dtAssembleCmd = &cobra.Command{
	Use:   "dtAssemble",
	Short: "helps assembling multiple DT project sboms into a final sbom",
	Long: `The assemble command will help assembling sboms into a final sbom.

Basic Example:
    $ sbomasm dtAssemble -u "http://localhost:8080/" -k "odt_gwiwooi29i1N5Hewkkddkkeiwi3ii" -n "mega-app" -v "1.0.0" -t "application" -o finalsbom.json 11903ba9-a585-4dfb-9a0c-f348345a5473 34103ba2-rt63-2fga-3a8b-t625261g6262
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please provide at least one sbom file to assemble")
		}

		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logger.InitDebugLogger()
		} else {
			logger.InitProdLogger()
		}

		ctx := logger.WithLogger(context.Background())

		dtAssembleParams, err := extractDTArgs(cmd, args)
		if err != nil {
			return err
		}

		dtAssembleParams.Ctx = &ctx

		// retrieve Input Files
		dtassemble.PopulateInputField(ctx, dtAssembleParams)
		fmt.Println("dtAssembleParams.Input: ", dtAssembleParams.Input)

		assembleParams, err := extractArgsFromDTAssembleToAssemble(dtAssembleParams)
		if err != nil {
			return err
		}
		assembleParams.Ctx = &ctx

		config, err := assemble.PopulateConfig(assembleParams)
		if err != nil {
			fmt.Println("Error populating config:", err)
		}
		return assemble.Assemble(config)
	},
}

func extractArgsFromDTAssembleToAssemble(dtAssembleParams *dtassemble.Params) (*assemble.Params, error) {
	aParams := assemble.NewParams()

	aParams.Output = dtAssembleParams.Output

	aParams.Name = dtAssembleParams.Name
	aParams.Version = dtAssembleParams.Version
	aParams.Type = dtAssembleParams.Type

	aParams.FlatMerge = dtAssembleParams.FlatMerge
	aParams.HierMerge = dtAssembleParams.HierMerge
	aParams.AssemblyMerge = dtAssembleParams.AssemblyMerge

	aParams.Xml = dtAssembleParams.Xml
	aParams.Json = dtAssembleParams.Json

	aParams.OutputSpecVersion = dtAssembleParams.OutputSpecVersion

	aParams.OutputSpec = dtAssembleParams.OutputSpec

	aParams.Input = dtAssembleParams.Input

	return aParams, nil
}

func init() {
	rootCmd.AddCommand(dtAssembleCmd)
	dtAssembleCmd.Flags().StringP("url", "u", "", "dependency track url https://localhost:8080/")
	dtAssembleCmd.Flags().StringP("api-key", "k", "", "dependency track api key, requires VIEW_PORTFOLIO for scoring and PORTFOLIO_MANAGEMENT for tagging")
	dtAssembleCmd.MarkFlagsRequiredTogether("url", "api-key")

	dtAssembleCmd.Flags().StringP("output", "o", "", "path to assembled sbom, defaults to stdout")

	dtAssembleCmd.Flags().StringP("name", "n", "", "name of the assembled sbom")
	dtAssembleCmd.Flags().StringP("version", "v", "", "version of the assembled sbom")
	dtAssembleCmd.Flags().StringP("type", "t", "", "product type of the assembled sbom (application, framework, library, container, device, firmware)")
	dtAssembleCmd.MarkFlagsRequiredTogether("name", "version", "type")

	dtAssembleCmd.Flags().BoolP("flatMerge", "f", false, "flat merge")
	dtAssembleCmd.Flags().BoolP("hierMerge", "m", false, "hierarchical merge")
	dtAssembleCmd.Flags().BoolP("assemblyMerge", "a", false, "assembly merge")
	dtAssembleCmd.MarkFlagsMutuallyExclusive("flatMerge", "hierMerge", "assemblyMerge")

	dtAssembleCmd.Flags().BoolP("outputSpecCdx", "g", true, "output in cdx format")
	dtAssembleCmd.Flags().BoolP("outputSpecSpdx", "s", false, "output in spdx format")
	dtAssembleCmd.MarkFlagsMutuallyExclusive("outputSpecCdx", "outputSpecSpdx")

	dtAssembleCmd.Flags().StringP("outputSpecVersion", "e", "", "spec version of the output sbom")

	dtAssembleCmd.Flags().BoolP("xml", "x", false, "output in xml format")
	dtAssembleCmd.Flags().BoolP("json", "j", true, "output in json format")
	dtAssembleCmd.MarkFlagsMutuallyExclusive("xml", "json")
}

func extractDTArgs(cmd *cobra.Command, args []string) (*dtassemble.Params, error) {
	aParams := dtassemble.NewParams()

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return nil, err
	}

	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return nil, err
	}
	aParams.Url = url
	aParams.ApiKey = apiKey

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}
	aParams.Output = output

	name, _ := cmd.Flags().GetString("name")
	version, _ := cmd.Flags().GetString("version")
	typeValue, _ := cmd.Flags().GetString("type")

	aParams.Name = name
	aParams.Version = version
	aParams.Type = typeValue

	flatMerge, _ := cmd.Flags().GetBool("flatMerge")
	hierMerge, _ := cmd.Flags().GetBool("hierMerge")
	assemblyMerge, _ := cmd.Flags().GetBool("assemblyMerge")

	aParams.FlatMerge = flatMerge
	aParams.HierMerge = hierMerge
	aParams.AssemblyMerge = assemblyMerge

	xml, _ := cmd.Flags().GetBool("xml")
	json, _ := cmd.Flags().GetBool("json")

	aParams.Xml = xml
	aParams.Json = json

	if aParams.Xml {
		aParams.Json = false
	}

	specVersion, _ := cmd.Flags().GetString("outputSpecVersion")
	aParams.OutputSpecVersion = specVersion

	cdx, _ := cmd.Flags().GetBool("outputSpecCdx")

	if cdx {
		aParams.OutputSpec = "cyclonedx"
	} else {
		aParams.OutputSpec = "spdx"
	}

	fmt.Println("args: ", args)
	for _, arg := range args {
		fmt.Println("arg: ", arg)
		argID, err := uuid.Parse(arg)
		fmt.Println("argID: ", argID)

		if err != nil {
			return nil, err
		}
		aParams.ProjectIds = append(aParams.ProjectIds, argID)
	}
	return aParams, nil
}
