// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/codelingo/hub/util/mdgen/dataStruct"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// genTenetsCmd represents the genTenets command
var genTenetsCmd = &cobra.Command{
	Use:   "genTenets",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		tmplPath := "template/tenet.md"
		//outPath := args[0]

		tenetsRootPath := os.Getenv("GOPATH") + "/src/github.com/codelingo/hub/tenets"
		result, err := parseTenetsDir(tenetsRootPath)
		if err != nil {
			panic(err)
		}
		for ownerKey, owner := range result.Owners {
			for bundleKey, bundle := range owner.Bundles {
				for tenetKey, _ := range bundle.Tenets {
					data := dataStruct.Data{
						Owner:  ownerKey,
						Bundle: bundleKey,
						Tenet:  tenetKey,
					}
					outPath := filepath.Join(tenetsRootPath, ownerKey, bundleKey, tenetKey, "README.md")

					err := writeFile(tmplPath, outPath, data)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genTenetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genTenetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genTenetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseTenetsDir(dirPath string) (dataStruct.HubTenets, error) {
	var hubTenets dataStruct.HubTenets
	if filepath.Base(dirPath) != "tenets" {
		return hubTenets, errors.New("Please select the correct tenets folder i.e. <...>/hub/tenets")
	}
	owners, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return hubTenets, errors.Trace(err)
	}

	hubTenets.Owners = make(map[string]dataStruct.TenetsOwner)
	for _, owner := range owners {
		if owner.IsDir() {
			hubTenets.Owners[owner.Name()] = dataStruct.TenetsOwner{Name: owner.Name()}
			ownerPath := filepath.Join(dirPath, owner.Name())
			bundles, err := ioutil.ReadDir(ownerPath)
			if err != nil {
				return hubTenets, errors.Trace(err)
			}
			tenetsOwner := hubTenets.Owners[owner.Name()]
			tenetsOwner.Bundles = make(map[string]dataStruct.Bundle)

			for _, bundle := range bundles {
				if bundle.IsDir() {
					tenetsOwner.Bundles[bundle.Name()] = dataStruct.Bundle{Name: bundle.Name()}
					hubTenets.Owners[owner.Name()] = tenetsOwner

					tenetPath := filepath.Join(ownerPath, bundle.Name())
					tenets, err := ioutil.ReadDir(tenetPath)
					if err != nil {
						return hubTenets, errors.Trace(err)
					}
					tmpBundle := hubTenets.Owners[owner.Name()].Bundles[bundle.Name()]
					tmpBundle.Tenets = make(map[string]dataStruct.Tenet)

					for _, tenet := range tenets {
						if tenet.IsDir() {
							tmpBundle.Tenets[tenet.Name()] = dataStruct.Tenet{Name: tenet.Name()}
							hubTenets.Owners[owner.Name()].Bundles[bundle.Name()] = tmpBundle

							hubTenets.Owners[owner.Name()].Bundles[bundle.Name()].Tenets[tenet.Name()] = dataStruct.Tenet{Name: tenet.Name()}
						}
					}
				}
			}
		}
	}

	return hubTenets, nil
}
