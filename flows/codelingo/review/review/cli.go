package review

import (
	"context"
	"fmt"
	"os"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/commands/verify"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/app/util/common/config"
	"github.com/codelingo/lingo/vcs"
	"github.com/codelingo/rpc/flow"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

const (
	vcsGit string = "git"
	vcsP4  string = "perforce"
)

var CLIApp = &flowutil.CLIApp{
	App: cli.App{

		Name:    "review",
		Usage:   "Review code following tenets in codelingo.yaml.",
		Version: "0.0.0",

		Flags: []cli.Flag{
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
				Usage: "Review a given directory.",
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
		// "$ lingo review" will review any unstaged changes from pwd down.
		// "$ lingo review [<filename>]" will review any unstaged changes in the named files.
		// "$ lingo review --all [<filename>]" will review all code in the named files.
	},
	Request: reviewAction,
}

func reviewRequire() error {
	reqs := []verify.Require{verify.VCSRq, verify.HomeRq, verify.AuthRq, verify.ConfigRq, verify.VersionRq}
	for _, req := range reqs {
		err := req.Verify()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func reviewAction(cliCtx *cli.Context) (chan proto.Message, chan error, func(), error) {
	err := reviewRequire()
	if err != nil {
		return nil, nil, nil, err
	}

	defer util.Logger.Sync()
	if cliCtx.IsSet("debug") {
		util.Logger.Debugw("reviewCMD called")
	}
	dir := cliCtx.String("directory")
	if dir != "" {
		if err := os.Chdir(dir); err != nil {
			return nil, nil, nil, err
		}
	}

	dotlingo, err := ReadDotLingo(cliCtx)
	if err != nil {
		return nil, nil, nil, err

	}
	vcsType, repo, err := vcs.New()
	if err != nil {
		return nil, nil, nil, err

	}

	// TODO: replace this system with nfs-like communication.
	fmt.Println("Syncing your repo...")
	if err = vcs.SyncRepo(vcsType, repo); err != nil {
		return nil, nil, nil, err

	}

	owner, name, err := repo.OwnerAndNameFromRemote()
	if err != nil {
		return nil, nil, nil, err

	}

	sha, err := repo.CurrentCommitId()
	if err != nil {
		if flowutil.NoCommitErr(err) {
			return nil, nil, nil, errors.New(flowutil.NoCommitErrMsg)
		}

		return nil, nil, nil, err

	}

	patches, err := repo.Patches()
	if err != nil {
		return nil, nil, nil, err

	}
	var patchesSize int64
	for _, patch := range patches {
		patchesSize += int64(len([]byte(patch)))
	}
	if patchesSize >= 1024*1024 { // >= 1MB; default max GRPC msg size accepted by the servers is 4MB.
		util.UserFacingWarning("Warning: large diffs can be error prone. You may need to commit your changes.")
	}

	workingDir, err := repo.WorkingDir()
	if err != nil {
		return nil, nil, nil, err

	}

	cfg, err := config.Platform()
	if err != nil {
		return nil, nil, nil, err

	}
	vcsTypeStr, err := vcs.TypeToString(vcsType)
	if err != nil {
		return nil, nil, nil, err

	}

	ctx, cancel := util.UserCancelContext(context.Background())
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
			return nil, nil, nil, err

		}
		hostname, err := cfg.GitRemoteName()
		if err != nil {
			return nil, nil, nil, err

		}

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Owner{owner}
	case vcsP4:
		addr, err := cfg.P4ServerAddr()
		if err != nil {
			return nil, nil, nil, err

		}
		hostname, err := cfg.P4RemoteName()
		if err != nil {
			return nil, nil, nil, err

		}
		depot, err := cfg.P4RemoteDepotName()
		if err != nil {
			return nil, nil, nil, err

		}
		name = owner + "/" + name

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Depot{depot}
		req.Repo = name
	default:
		return nil, nil, nil, errors.Errorf("Invalid VCS '%s'", vcsTypeStr)
	}

	fmt.Println("Running review flow...")
	resultc, errc, err := RequestReview(ctx, req)
	return resultc, errc, cancel, errors.Trace(err)
}
