package review

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
			cli.StringFlag{
				Name:  "repo",
				Usage: "Review a repo directly, e.g. github.com/some/repo",
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

func alwaysRequire() error {
	reqs := []verify.Require{verify.HomeRq, verify.ConfigRq, verify.VersionRq}
	for _, req := range reqs {
		if err := req.Verify(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func noRepoRequire() error {
	reqs := []verify.Require{verify.VCSRq, verify.AuthRq}
	for _, req := range reqs {
		if err := req.Verify(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func reviewAction(cliCtx *cli.Context) (chan proto.Message, <-chan *flowutil.UserVar, chan error, func(), error) {
	if err := alwaysRequire(); err != nil {
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

	ctx, cancel := util.UserCancelContext(context.Background())

	if repo := cliCtx.String("repo"); repo != "" {
		parts := strings.Split(repo, "/")
		if len(parts) != 3 {
			msg := "cannot parse repo; expected e.g. github.com/someuser/somerepo"
			return nil, nil, nil, nil, errors.New(msg)
		}
		req := &flow.ReviewRequest{
			Vcs:          "git",
			Host:         parts[0],
			OwnerOrDepot: &flow.ReviewRequest_Owner{parts[1]},
			Repo:         parts[2],
			Dotlingo:     dotlingo,
		}

		fmt.Println("Running review flow...")
		resultc, userVarc, errc, err := RequestReview(ctx, req, insecure)
		return resultc, userVarc, errc, cancel, errors.Trace(err)
	}

	if err := noRepoRequire(); err != nil {
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

	path, err := findPath()
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)

	}

	req := &flow.ReviewRequest{
		Repo:     name,
		Sha:      sha,
		Patches:  patches,
		Vcs:      vcsTypeStr,
		Dir:      workingDir,
		Dotlingo: dotlingo,
		Path:     path,
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

// `path` is the relative path from $GOPATH/src to the repo root. It's required for package resolution
// It is empty for non-Go repos
func findPath() (string, error) {

	// Find the root of repo
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Trace(err)
	}
	repoRoot := strings.TrimSpace(string(output))

	// Check if the repo is a go project
	goProject, err := isGoProject(repoRoot)
	if err != nil {
		return "", errors.Trace(err)
	}

	// Check if GOPATH is set
	goPath := os.Getenv("GOPATH")

	// Check if repo is in GOPATH
	if goPath != "" && goProject {
		if strings.HasPrefix(repoRoot, goPath+"/src/") {
			return strings.TrimPrefix(repoRoot, goPath+"/src/"), nil
		}
	}
	return "", nil
}

func isGoProject(currentDir string) (bool, error) {
	items, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return false, errors.Trace(err)
	}
	for _, item := range items {
		if !item.IsDir() && filepath.Ext(item.Name()) == ".go" {
			return true, nil
		} else if item.IsDir() && !strings.HasPrefix(item.Name(), ".") {
			isGo, err := isGoProject(filepath.Join(currentDir, item.Name()))
			if err != nil {
				return false, errors.Trace(err)
			}
			if isGo {
				return true, nil
			}
		}
	}
	return false, nil
}
