package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v1"
)

func main() {
	actualRules := getActualRules()
	statedRules := getStatedRules()
	missingActualRules := []string{}
	missingStatedRules := []string{}

	for rule := range actualRules {
		if _, ok := statedRules[rule]; !ok {
			missingStatedRules = append(missingStatedRules, rule)
		}
	}

	for rule := range statedRules {
		if _, ok := actualRules[rule]; !ok {
			missingActualRules = append(missingActualRules, rule)
		}
	}

	fmt.Println("Missing Stated Rules: ", missingStatedRules)
	fmt.Println("Missing Actual Rules: ", missingActualRules)
}

func getStatedRules() map[string]bool {
	rules := map[string]bool{}
	gopath := os.Getenv("GOPATH")
	registryPath := path.Join(gopath, "src/github.com/codelingo/codelingo/registry/tenets.yaml")

	registryContents, err := ioutil.ReadFile(registryPath)
	if err != nil {
		panic(err.Error())
	}

	r := map[interface{}]map[interface{}]map[interface{}]map[interface{}]interface{}{}
	yaml.Unmarshal(registryContents, &r)
	for o, bundles := range r {
		owner := o.(string)
		for b, contents := range bundles {
			bundle := b.(string)
			for key, ruleList := range contents {
				k := key.(string)
				if k != "tenets" {
					continue
				}
				for r := range ruleList {
					rule := r.(string)
					rules[path.Join(owner, bundle, rule)] = true
				}
			}
		}
	}
	return rules
}

func getActualRules() map[string]bool {
	gopath := os.Getenv("GOPATH")
	rulesPath := path.Join(gopath, "src/github.com/codelingo/codelingo/tenets")

	rules := map[string]bool{}

	owners, err := ioutil.ReadDir(rulesPath)
	if err != nil {
		panic(err.Error())
	}

	for _, owner := range owners {
		if owner.IsDir() {
			bundles, err := ioutil.ReadDir(path.Join(rulesPath, owner.Name()))
			if err != nil {
				panic(err.Error())
			}

			for _, bundle := range bundles {
				if bundle.IsDir() {
					ruleDirs, err := ioutil.ReadDir(path.Join(rulesPath, owner.Name(), bundle.Name()))
					if err != nil {
						panic(err.Error())
					}

					for _, rule := range ruleDirs {
						if rule.IsDir() {
							rules[path.Join(owner.Name(), bundle.Name(), rule.Name())] = true
						}
					}
				}
			}

		}
	}

	return rules
}
