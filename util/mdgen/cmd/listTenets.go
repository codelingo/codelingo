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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/codelingo/hub/util/mdgen/dataStruct"
	"github.com/spf13/cobra"
)

// listTenetsCmd represents the listTenets command
var listTenetsCmd = &cobra.Command{
	Use:   "listTenets",
	Short: "Lists out all available tenets in hub",
	Long: `Generates a list of all hub tenets, which can be used
statically in the website etc. There are no file generating side
effects and output is json of the form:

{
	name: 'codelingo/hub - nil_only_functions',
	repo: 'github.com/codelingo/hub',
	dir: '/',
	tenet: 'codelingo/go/nil_only_functions'
},
	`,
	Run: func(cmd *cobra.Command, args []string) {
		tenetsRootPath := os.Getenv("GOPATH") + "/src/github.com/codelingo/hub/tenets"
		result, err := parseTenetsDir(tenetsRootPath)
		if err != nil {
			panic(err)
		}

		tenetList := buildTenetList(result)
		sort.Sort(byUserPriority(tenetList))

		output, err := json.MarshalIndent(tenetList, "  ", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(output))

	},
}

type byUserPriority []*TenetDesc

func (list byUserPriority) Len() int {
	return len(list)
}
func (list byUserPriority) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list byUserPriority) Less(i, j int) bool {
	// Tenets are already grouped by owner due to traversal, and we priorities codelingo tenets
	priorityI := priorityByOwner(list[i].owner)
	priorityJ := priorityByOwner(list[j].owner)
	return priorityI > priorityJ
}

func priorityByOwner(owner string) int {
	if priority, ok := priorities[owner]; ok {
		return priority
	}
	return 0
}

var priorities = map[string]int{
	"codelingo": 1,
	"west":      -1,
}

// TenetDesc is a description of a tenet that can be used by the hub
type TenetDesc struct {
	Name  string `json:"name"`
	Repo  string `json:"repo"`
	Dir   string `json:"dir"`
	Tenet string `json:"tenet"`
	owner string
}

func buildTenetList(hubtenets dataStruct.HubTenets) []*TenetDesc {
	tenetList := []*TenetDesc{}

	for ownerKey, owner := range hubtenets.Owners {
		for bundleKey, bundle := range owner.Bundles {
			for tenetKey := range bundle.Tenets {
				tenetName := filepath.Join(ownerKey, bundleKey, tenetKey)

				tenetList = append(tenetList, &TenetDesc{
					Name:  "codelingo/hub - " + tenetKey,
					Repo:  "github.com/codelingo/hub",
					Dir:   "tenets/" + tenetName,
					Tenet: tenetName,
					owner: ownerKey,
				})
			}
		}
	}

	return tenetList
}

func init() {
	rootCmd.AddCommand(listTenetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listTenetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listTenetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
