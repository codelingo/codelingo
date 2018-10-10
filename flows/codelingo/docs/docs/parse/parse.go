package parse

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/codelingo/clql/dotlingo"
	"github.com/juju/errors"
)

// TODO: hang Parse off a worker that user that conforms to the juju worker interface.
// TODO: parse currently just calls the dotlingo parser, but eventually it should call out to the
// dotlingo/CLQL lexicon.

// Parse is responsible for all dotlingo parsing in the bot flow layer.
func Parse(dl string) ([]*dotlingo.Dotlingo, error) {
	if dl == "" {
		return nil, errors.New("cannot parse codelingo, string is empty")
	}
	queries, err := dotlingo.Parse(dl)
	if err != nil {
		return nil, errors.Annotate(err, "failed to parse given lingo")
	}

	allQueries, err := resolveImports(queries)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return allQueries, nil
}

// resolveImports retrieves all bundles and tenets imported in a given codelingo.yaml file, parses them,
// and returns them in a list including the original one.
// TODO: resolve recursively
func resolveImports(parentQuery *dotlingo.Dotlingo) ([]*dotlingo.Dotlingo, error) {
	allQueries := []*dotlingo.Dotlingo{parentQuery}

	for _, genericImp := range parentQuery.Imports {
		switch imp := genericImp.ImportType.(type) {
		case *dotlingo.Import_Tenet:
			path := imp.Tenet
			importedDotLingo, err := requestDotLingo(path.Owner, path.Bundle, path.Name)
			if err != nil {
				return nil, errors.Trace(err)
			}

			importedQuery, err := dotlingo.Parse(importedDotLingo)
			if err != nil {
				return nil, errors.Annotatef(err, "parse failed for imported tenet %s/%s/%s", path.Owner, path.Bundle, path.Name)
			}

			allQueries = append(allQueries, importedQuery)
		case *dotlingo.Import_Bundle:
			path := imp.Bundle
			importedDls, err := requestBundle(path.Owner, path.Name)
			if err != nil {
				return nil, errors.Trace(err)
			}

			for dlname, dl := range importedDls {
				query, err := dotlingo.Parse(dl)
				if err != nil {
					return nil, errors.Annotatef(err, "parse failed for imported tenet %s/%s/%s", path.Owner, path.Name, dlname)
				}

				allQueries = append(allQueries, query)
			}
		default:
			return nil, errors.Errorf("unexpected import type %v", imp)
		}
	}

	return allQueries, nil
}

const github404 = "404: Not Found\n"

const tenetURL = "https://raw.githubusercontent.com/codelingo/codelingo/master/tenets/%s/%s/%s/codelingo.yaml"

// Request a DotLingo from the github.com/codelingo/codelingo repo
func requestDotLingo(owner, bundle, name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(tenetURL, owner, bundle, name))
	if err != nil {
		return "", errors.Trace(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	if string(body) == github404 {
		return "", errors.Errorf("%s there's no tenet called \"%s\" at github.com/codelingo/codelingo/tenets/%s/%s", github404, name, owner, bundle)
	}

	return string(body), nil
}

const bundleURL = "https://raw.githubusercontent.com/codelingo/codelingo/master/tenets/%s/%s/lingo_bundle.yaml"

type rawLingoBundle struct {
	Tenets []string `yaml:"tenets"`
	Tags   []string `yaml:"tags"`
}

// Requests a Tenet Bundle from the github.com/codelingo/codelingo repo
func requestBundle(owner, bundleName string) (map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf(bundleURL, owner, bundleName))
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if string(body) == github404 {
		return nil, errors.Errorf("%s there's no lingo_bundle called \"%s\" at github.com/codelingo/codelingo/tenets/%s", github404, bundleName, owner)
	}

	var rawBun *rawLingoBundle
	err = yaml.Unmarshal([]byte(body), &rawBun)
	if err != nil {
		return nil, errors.Trace(err)
	}

	dls := map[string]string{}
	for _, tenetName := range rawBun.Tenets {
		// TODO: optimize with concurrency
		newDl, err := requestDotLingo(owner, bundleName, tenetName)
		if err != nil {
			return nil, errors.Trace(err)
		}

		dls[tenetName] = newDl
	}

	return dls, nil
}
