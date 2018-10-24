package rewrite

import (
	"fmt"
	"os"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	"github.com/urfave/cli"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/commands/verify"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/app/util/common/config"
	"github.com/codelingo/lingo/vcs"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
)

const (
	vcsGit string = "git"
	vcsP4  string = "perforce"
)

var CLIApp = &flowutil.CLIApp{
	Tagline: "Automated Code Fixes",
	App: cli.App{
		Name:    "rewrite",
		Usage:   "The Rewrite Flow rewrites sections of source code matching the query pattern in the Tenets it's run over.",
		Version: "0.0.0",
		// Action:  action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  util.OutputFlg.String(),
				Usage: "File to save found rewrites to.",
			},
			cli.StringFlag{
				Name:  util.FormatFlg.String(),
				Value: "json-pretty",
				Usage: "How to format the results. Possible values are: json, json-pretty.",
			},
			cli.BoolFlag{
				Name:  util.KeepAllFlg.String(),
				Usage: "Keep all rewrites and don't be prompted to confirm each.",
			},
			cli.StringFlag{
				Name:  "dir, d",
				Usage: "Modify a given directory.",
			},
			cli.BoolFlag{
				Name:  "debug",
				Usage: "Display debug messages",
			},
		},
	},
	Request: rewriteAction,
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

func rewriteAction(cliCtx *cli.Context) (chan proto.Message, chan error, func(), error) {
	err := rewriteRequire()
	if err != nil {
		util.FatalOSErr(err)
		return nil, nil, nil, err
	}

	defer util.Logger.Sync()
	if cliCtx.IsSet("debug") {
		util.Logger.Debugw("CliCMD called")
	}
	dir := cliCtx.String("directory")
	if dir != "" {
		if err := os.Chdir(dir); err != nil {
			return nil, nil, nil, errors.Trace(err)
		}
	}

	dotlingo, err := flowutil.ReadDotLingo(cliCtx)
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}
	vcsType, repo, err := vcs.New()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	// TODO: replace this system with nfs-like communication.
	fmt.Println("Syncing your repo...")
	if err = vcs.SyncRepo(vcsType, repo); err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	owner, name, err := repo.OwnerAndNameFromRemote()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	sha, err := repo.CurrentCommitId()
	if err != nil {
		if flowutil.NoCommitErr(err) {
			return nil, nil, nil, errors.New(flowutil.NoCommitErrMsg)
		}

		return nil, nil, nil, errors.Trace(err)
	}

	patches, err := repo.Patches()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	workingDir, err := repo.WorkingDir()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	cfg, err := config.Platform()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}
	vcsTypeStr, err := vcs.TypeToString(vcsType)
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	req := &rewriterpc.Request{
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
			return nil, nil, nil, errors.Trace(err)
		}
		hostname, err := cfg.GitRemoteName()
		if err != nil {
			return nil, nil, nil, errors.Trace(err)
		}

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &rewriterpc.Request_Owner{owner}
	case vcsP4:
		addr, err := cfg.P4ServerAddr()
		if err != nil {
			return nil, nil, nil, errors.Trace(err)
		}
		hostname, err := cfg.P4RemoteName()
		if err != nil {
			return nil, nil, nil, errors.Trace(err)
		}
		depot, err := cfg.P4RemoteDepotName()
		if err != nil {
			return nil, nil, nil, errors.Trace(err)
		}
		name = owner + "/" + name

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &rewriterpc.Request_Depot{depot}
		req.Repo = name
	default:
		return nil, nil, nil, errors.Errorf("Invalid VCS '%s'", vcsTypeStr)
	}

	fmt.Println("Running rewrite flow...")

	// proto.RegisterType((*rewriterpc.Hunk)(nil), "rpc.Hunk")
	return flowutil.RunFlow("rewrite", req, func() proto.Message { return &rewriterpc.Hunk{} })
}
