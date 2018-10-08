package main

import (
	"context"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite"

	"github.com/codelingo/lingo/app/commands/verify"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/app/util/common/config"
	"github.com/codelingo/lingo/vcs"
	"github.com/codelingo/rpc/flow"
	flowutil "github.com/codelingo/sdk/flow"
	"github.com/juju/errors"
)

const (
	vcsGit string = "git"
	vcsP4  string = "perforce"
)

var rewriteCmd = cli.Command{
	Name:  "rewrite",
	Usage: "Modify code following tenets in codelingo.yaml.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  util.LingoFile.String(),
			Usage: "A codelingo.yaml file to perform the review with. If the flag is not set, codelingo.yaml files are read from the branch being reviewed.",
		},
		cli.StringFlag{
			Name:  util.DiffFlg.String(),
			Usage: "Review only unstaged changes in the working tree.",
		},
		cli.StringFlag{
			Name:  util.OutputFlg.String(),
			Usage: "File to save found issues to.",
		},
		cli.StringFlag{
			Name:  util.FormatFlg.String(),
			Value: "json-pretty",
			Usage: "How to format the found issues. Possible values are: json, json-pretty.",
		},
		cli.BoolFlag{
			Name:  util.KeepAllFlg.String(),
			Usage: "Keep all issues and don't be prompted to confirm each issue.",
		},
		cli.StringFlag{
			Name:  util.DirectoryFlg.String(),
			Usage: "Modify a given directory.",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Display debug messages",
		},
		// cli.BoolFlag{
		// 	Name:  "all",
		// 	Usage: "review all files under all directories from pwd down",
		// },
	},
	Description: `
"$ lingo review" will review all code from pwd down.
"$ lingo review <filename>" will only review named file.
`[1:],
	// "$ lingo review" will review any unstaged changes from pwd down.
	// "$ lingo review [<filename>]" will review any unstaged changes in the named files.
	// "$ lingo review --all [<filename>]" will review all code in the named files.
	Action: rewriteAction,
}

func main() {
	if err := flowutil.Run(rewriteCmd); err != nil {
		flowutil.HandleErr(err)
	}
}

func rewriteAction(ctx *cli.Context) {
	err := rewriteRequire()
	if err != nil {
		util.FatalOSErr(err)
		return
	}

	msg, err := rewriteCMD(ctx)
	if err != nil {
		if ctx.IsSet("debug") {
			// Debugging
			util.Logger.Debugw("rewriteAction", "err_stack", errors.ErrorStack(err))
		}
		util.FatalOSErr(err)
		return
	}

	fmt.Println(msg)
}

func rewriteRequire() error {
	reqs := []verify.Require{verify.VCSRq, verify.HomeRq, verify.AuthRq, verify.ConfigRq, verify.VersionRq}
	for _, req := range reqs {
		err := req.Verify()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func rewriteCMD(cliCtx *cli.Context) (string, error) {
	defer util.Logger.Sync()
	if cliCtx.IsSet("debug") {
		util.Logger.Debugw("rewriteCMD called")
	}
	dir := cliCtx.String("directory")
	if dir != "" {
		if err := os.Chdir(dir); err != nil {
			return "", errors.Trace(err)
		}
	}

	dotlingo, err := rewrite.ReadDotLingo(cliCtx)
	if err != nil {
		return "", errors.Trace(err)
	}
	vcsType, repo, err := vcs.New()
	if err != nil {
		return "", errors.Trace(err)
	}

	// TODO: replace this system with nfs-like communication.
	fmt.Println("Syncing your repo...")
	if err = vcs.SyncRepo(vcsType, repo); err != nil {
		return "", errors.Trace(err)
	}

	owner, name, err := repo.OwnerAndNameFromRemote()
	if err != nil {
		return "", errors.Trace(err)
	}

	sha, err := repo.CurrentCommitId()
	if err != nil {
		if flowutil.NoCommitErr(err) {
			return "", errors.New(flowutil.NoCommitErrMsg)
		}

		return "", errors.Trace(err)
	}

	patches, err := repo.Patches()
	if err != nil {
		return "", errors.Trace(err)
	}

	workingDir, err := repo.WorkingDir()
	if err != nil {
		return "", errors.Trace(err)
	}

	cfg, err := config.Platform()
	if err != nil {
		return "", errors.Trace(err)
	}
	vcsTypeStr, err := vcs.TypeToString(vcsType)
	if err != nil {
		return "", errors.Trace(err)
	}

	ctx, cancel := util.UserCancelContext(context.Background())
	issuec := make(chan *flow.Issue)
	errorc := make(chan error)

	req := &flow.ReviewRequest{
		Repo:     name,
		Sha:      sha,
		Patches:  patches,
		Vcs:      vcsTypeStr,
		Dir:      workingDir,
		Dotlingo: dotlingo,
	}
	switch vcsTypeStr {
	case vcsGit:
		addr, err := cfg.GitServerAddr()
		if err != nil {
			return "", errors.Trace(err)
		}
		hostname, err := cfg.GitRemoteName()
		if err != nil {
			return "", errors.Trace(err)
		}

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Owner{owner}
	case vcsP4:
		addr, err := cfg.P4ServerAddr()
		if err != nil {
			return "", errors.Trace(err)
		}
		hostname, err := cfg.P4RemoteName()
		if err != nil {
			return "", errors.Trace(err)
		}
		depot, err := cfg.P4RemoteDepotName()
		if err != nil {
			return "", errors.Trace(err)
		}
		name = owner + "/" + name

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Depot{depot}
		req.Repo = name
	default:
		return "", errors.Errorf("Invalid VCS '%s'", vcsTypeStr)
	}

	fmt.Println("Running rewrite flow...")
	issuec, errorc, err = rewrite.RequestReview(ctx, req)
	if err != nil {
		return "", errors.Trace(err)
	}

	issues, err := rewrite.ConfirmIssues(cancel, issuec, errorc, cliCtx.Bool("keep-all"), cliCtx.String("output"))
	if err != nil {
		return "", errors.Trace(err)
	}

	if len(issues) == 0 {
		return fmt.Sprintf("Done! No issues found.\n"), nil
	}

	// Remove dicarded issues from report
	var keptIssues []*rewrite.SRCHunk
	for _, issue := range issues {
		if !issue.Discard {
			keptIssues = append(keptIssues, issue)
		}
	}
	return "", errors.Trace(rewrite.Write(keptIssues))
}
