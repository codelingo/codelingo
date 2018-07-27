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
	"fmt"
	"github.com/codelingo/lingo/service"
	"github.com/juju/errors"
	"github.com/marcinwyszynski/directory_tree"
	"github.com/spf13/cobra"
	"os"
)

// genLexiconCmd represents the genLexicon command
var genLexiconCmd = &cobra.Command{
	Use:   "genLexicon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		lexs, err := listLexs()
		if err != nil {
			panic(err)
		}

		for _, lex := range lexs {
			if err := writeLexMD(lex); err != nil {
				panic(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(genLexiconCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genLexiconCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genLexiconCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type lexInfo struct {
	Typ, Owner, Name, OutPath string
	Facts                     map[string][]string
}

func listLexs() ([]*lexInfo, error) {

	lexsPath := os.Getenv("GOPATH") + "/src/github.com/codelingo/hub/lexicons"

	tree, err := directory_tree.NewTree(lexsPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var lexInfos []*lexInfo

	for _, typNode := range tree.Children {
		if !typNode.Info.IsDir {
			continue
		}
		typName := typNode.Info.Name
		for _, ownerNode := range typNode.Children {
			if !ownerNode.Info.IsDir {
				continue
			}
			ownerName := ownerNode.Info.Name
			for _, lexNode := range ownerNode.Children {
				if !lexNode.Info.IsDir {
					continue
				}
				lexName := lexNode.Info.Name

				lexInfos = append(lexInfos, &lexInfo{
					Typ:     typName,
					Owner:   ownerName,
					Name:    lexName,
					OutPath: fmt.Sprintf("%s/%s/%s/%s", lexsPath, typName, ownerName, lexName),
				})
			}
		}
	}
	return lexInfos, nil
}

func writeLexMD(data *lexInfo) error {

	facts, err := listFacts(data.Owner, data.Name)
	if err != nil {
		// TODO(waigani) once list facts works for all lexicons, error here
		fmt.Printf("owner: %s, name: %s, error: %s", data.Owner, data.Name, err.Error())
	}

	outPath := data.OutPath + "/README.md"
	data.Facts = facts
	return writeFile(os.Getenv("GOPATH")+"/src/github.com/codelingo/hub/util/mdgen/template/lexicon.md", outPath, data)

}

func listFacts(owner, lexName string) (map[string][]string, error) {
	svc, err := service.New()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return svc.ListFacts(owner, lexName, "")
}
