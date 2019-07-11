package review

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
			cli.BoolFlag{
				Name:   "insecure",
				Hidden: true,
				Usage:  "Review without TLS",
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
	Tagline: "Automated Code Reviews",
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

func reviewAction(cliCtx *cli.Context) (chan proto.Message, <-chan *flowutil.UserVar, chan error, func(), error) {
	err := reviewRequire()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)
	}

	defer util.Logger.Sync()
	if cliCtx.IsSet("debug") {
		util.Logger.Debugw("reviewCMD called")
	}
	dir := cliCtx.String("directory")
	if dir != "" {
		if err := os.Chdir(dir); err != nil {
			return nil, nil, nil, nil, errors.Trace(err)
		}
	}

	insecure := cliCtx.IsSet("insecure")
	util.Logger.Debugf("insecure %t", insecure)

	dotlingo, err := ReadDotLingo(cliCtx)
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}
	vcsType, repo, err := vcs.New()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}

	// TODO: replace this system with nfs-like communication.
	fmt.Println("Syncing your repo...")
	if err = vcs.SyncRepo(vcsType, repo); err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}

	owner, name, err := repo.OwnerAndNameFromRemote()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}

	sha, err := repo.CurrentCommitId()
	if err != nil {
		if flowutil.NoCommitErr(err) {
			return nil, nil, nil, nil, errors.New(flowutil.NoCommitErrMsg)
		}

		return nil, nil, nil, nil, errors.Trace(err)

	}

	patches, err := repo.Patches()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

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
		return nil, nil, nil, nil, errors.Trace(err)

	}

	cfg, err := config.Platform()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}
	vcsTypeStr, err := vcs.TypeToString(vcsType)
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}
	// Build `path` for go projects as "reponame" for a local repo or "github.com/username/reponame" for a github repo
	// path variable is initially empty and will remain empty for none go repos
	path := ""

	// Find the current working directory, i.e. where user executed `lingo run review` command on the repo
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}
	// Check if the repo is a go project
	goProject, err := isGoProject(wd)
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)
	}

	// Check if GOPATH is set
	goPath := os.Getenv("GOPATH")

	if goPath != "" && goProject {
		if strings.Contains(wd, "github.com") { // Check if the working directory includes "github.com/username/reponame"
			swd := wd[strings.Index(wd, "github.com"):]
			parts := strings.Split(swd, "/")
			if len(parts) == 3 {
				path = swd
			}
		} else if strings.HasPrefix(wd, goPath+"/src/") { //Check if repo is local and lives in GOPATH/src/reponame
			swd := wd[strings.Index(wd, "src/"):]
			parts := strings.Split(swd, "/")
			if len(parts) == 2 {
				path = parts[1]
			}
		}
		// If none of the above, path remains empty
	}

	ctx, cancel := util.UserCancelContext(context.Background())

	var req *flow.ReviewRequest
	// If path is not empty, it is sent in ReviewRequest and used in go lexicon
	if path != "" {
		req = &flow.ReviewRequest{
			Repo:     name,
			Sha:      sha,
			Patches:  patches,
			Vcs:      vcsTypeStr,
			Dir:      workingDir,
			Dotlingo: dotlingo,
			Path:     path,
		}
	} else {
		req = &flow.ReviewRequest{
			Repo:     name,
			Sha:      sha,
			Patches:  patches,
			Vcs:      vcsTypeStr,
			Dir:      workingDir,
			Dotlingo: dotlingo,
		}
	}

	switch vcsTypeStr {
	case vcsGit:
		addr, err := cfg.GitServerAddr()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)

		}
		hostname, err := cfg.GitRemoteName()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)

		}

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Owner{owner}
	case vcsP4:
		addr, err := cfg.P4ServerAddr()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)

		}
		hostname, err := cfg.P4RemoteName()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)

		}
		depot, err := cfg.P4RemoteDepotName()
		if err != nil {
			return nil, nil, nil, nil, errors.Trace(err)

		}
		name = owner + "/" + name

		req.Host = addr
		req.Hostname = hostname
		req.OwnerOrDepot = &flow.ReviewRequest_Depot{depot}
		req.Repo = name
	default:
		return nil, nil, nil, nil, errors.Errorf("Invalid VCS '%s'", vcsTypeStr)
	}

	fmt.Println("Running review flow...")
	resultc, userVarc, errc, err := RequestReview(ctx, req, insecure)
	return resultc, userVarc, errc, cancel, errors.Trace(err)
}

func isGoProject(currentDir string) (bool, error) {
	isGo := false
	items, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.IsDir() {
			isGo, err = isGoProject(filepath.Join(currentDir, item.Name()))
			if err != nil {
				return false, err
			}
		} else if filepath.Ext(item.Name()) == ".go" {
			isGo = true
			break
		}
	}
	return isGo, nil
}
