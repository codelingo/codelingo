package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/commands/common/config"
	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
)

var reviewboard = cli.Command{

	Name:    "reviewboard",
	Aliases: []string{"rb"},
	Usage:   "post a review to a reviewboard service",
	Description: `
Connection details for the service will be read from $LINGO_HOME/services.yaml

The default the config under "reviewboard" will be read. To read another config block within services.yaml, use the --config flag.
`[1:],
	Flags: append(review.Flags,
		cli.StringFlag{
			// TODO(waigani) interactively set options for tenet.
			Name:  "config",
			Value: "reviewboard",
			Usage: "the name of the config block in services.yaml to use. It must be for a reviewboard service.",
		},
	),
	Action: rbAction,
}

func rbAction(ctx *cli.Context) {
	if err := rb(ctx); err != nil {
		common.OSErrf(err.Error())
		return
	}
	fmt.Println("Done! Review sent to Review Board.")
}

func rb(ctx *cli.Context) error {

	if len(ctx.Args()) == 0 {
		return errors.Errorf("The ReviewBoard service requires at least one arguement, the review ID")
	}

	// Make sure we have a config before posting
	serviceName := ctx.String("config")
	var cfgNeedsEdit bool
	rbCfg, err := config.Service(serviceName)
	if err != nil {
		fmt.Printf("no configuration found for %q\n", serviceName)
		cfgNeedsEdit = true
	} else if err := validateRbCfg(rbCfg); err != nil {
		fmt.Printf("%q configuration file is not valid: %v\n", serviceName, err)
		cfgNeedsEdit = true
	}
	if cfgNeedsEdit {
		// TODO(waigani) send and check a typed error
		var edit string
		fmt.Printf("Would you like to configure %s now? n/Y: ", serviceName)
		fmt.Scanln(&edit)

		switch edit {
		case "n", "N":
			return errors.New("aborted")
		case "y", "Y", "":
			return config.Edit(config.ServicesCfgFile, "vi") // TODO(waigani) use default editor
		default:
			return errors.Trace(err)
		}
		return nil
	}

	opts := review.Options{
		Files:      ctx.Args()[1:],
		Diff:       ctx.Bool("diff"),
		SaveToFile: ctx.String("save"),
		KeepAll:    ctx.Bool("keep-all"),
	}

	issues, err := review.Review(opts)
	if err != nil {
		return errors.Trace(err)
	}
	reviewID := ctx.Args()[0]

	fmt.Println("Posting review to reviewboard ...")

	return errors.Trace(postToRB(reviewID, rbCfg, issues))
}

func postToRB(reviewID string, rbCfg config.ServiceConfig, issues []*tenet.Issue) error {

	rbCfg["review-id"] = reviewID

	res := rbResult{
		rbCfg,
		issues,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return errors.Trace(err)
	}

	url := serviceEndpoint("reviewboard")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg := resp.Status
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Trace(err)
		}
		msg += "\n" + string(b)
		msg += "\n" + "if you are using the copr command, please use the --fetch-all flag to ensure the pull request's base branch is up-to-date"
		return errors.Errorf("failed to post to reviewboard service: %s", msg)
	}
	return nil
}

type rbResult struct {
	Config config.ServiceConfig `json:"config"`
	Issues []*tenet.Issue `json:"issues"`
}

func validateRbCfg(rbCfg config.ServiceConfig) error {
	if _, ok := rbCfg["domain"].(string); !ok {
		return errors.New("reviewboard domain not set in config")
	}
	if _, ok := rbCfg["token"].(string); !ok {
		return errors.New("reviewboard API token not set in config")
	}
	if _, ok := rbCfg["publish"].(string); !ok {
		rbCfg["publish"] = false
	}

	return nil
}
