package flow

import (
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/codelingo/lingo/app/util"
	"github.com/common-nighthawk/go-figure"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

type DecoratedResult struct {

	// Ctx are the flags and args parsed from the decorated fact
	Ctx *cli.Context

	// Payload is the struct returned form the flow server for the matched
	// fact.
	Payload proto.Message
}

type flowRunner struct {
	cliCtx           *cli.Context
	cliApp           *CLIApp
	decoratorApp     *DecoratorApp
	decoratedResultc chan *DecoratedResult
	waitc            chan struct{}
}

type CLIApp struct {
	cli.App
	Request func(*cli.Context) (chan proto.Message, chan error, func(), error)
}

type DecoratorApp struct {
	cli.App
	ConfirmDecorated func(*cli.Context, proto.Message) (bool, error)

	// Help info
	DecoratorUsage   string
	DecoratorExample string
}

func NewFlow(CLIApp *CLIApp, decoratorApp *DecoratorApp) *flowRunner {
	setBaseApp(CLIApp)

	fRunner := &flowRunner{
		cliApp:           CLIApp,
		decoratorApp:     decoratorApp,
		decoratedResultc: make(chan *DecoratedResult),
		waitc:            make(chan struct{}),
	}
	fRunner.setHelp()
	go func() {
		<-fRunner.waitc
		close(fRunner.decoratedResultc)
	}()

	//fRunner.cliApp.Action = fRunner.action
	return fRunner
}

func (f *flowRunner) safeCloseWaitc() {
	select {
	case <-f.waitc:
	default:
	}

	close(f.waitc)
}

func (f *flowRunner) setHelp() {
	defer f.safeCloseWaitc()

	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		figure.NewFigure("codelingo", "larry3d", false).Print()
		printHelp(w, CLIAPPHELPTMP, data)
		if f.decoratorApp != nil {
			printHelp(w, DECAPPHELPTMP, f.decoratorApp)
		}
		printHelp(w, INFOTMP, data)
	}
}

func setBaseApp(cliApp *CLIApp) {

	app := CLIApp{
		App:     *cli.NewApp(),
		Request: cliApp.Request,
	}

	// base settings
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lingo-file",
			Value: "codelingo.yaml",
			Usage: "the Tenet `FILE` to run rewrite over. If the flag is not set, codelingo.yaml files are read from the branch being rewritten.",
			//	Destination: &language,
		},
	}
	app.Compiled = time.Now()
	app.EnableBashCompletion = true

	// user settings
	app.Name = cliApp.Name
	app.Usage = cliApp.Usage
	app.Flags = append(app.Flags, cliApp.Flags...)
	if app.Action == nil {
		app.Action = cliApp.Action
	}
	app.Version = cliApp.Version

	*cliApp = app
}

func (f *flowRunner) Run() (_ chan *DecoratedResult, err error) {
	return f.decoratedResultc, f.cliApp.Run(os.Args)
}

func (f *flowRunner) action(ctx *cli.Context) {
	panic("action run")

	if err := f.command(ctx); err != nil {
		util.Logger.Debugw(errors.ErrorStack(err))
		util.FatalOSErr(err)
	}

}

func (f *flowRunner) command(ctx *cli.Context) (err error) {
	defer f.safeCloseWaitc()

	resultc, errc, cancel, err := f.RunCLI()
	if err != nil {
		return errors.Trace(err)
	}

	// If user is manually confirming results, set a long timeout.
	timeout := time.After(time.Hour * 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Wait()
		f.safeCloseWaitc()
	}()

l:
	for {
		select {
		case err, ok := <-errc:
			if !ok {
				errc = nil
				break
			}

			util.Logger.Debugf("Result error: %s", errors.ErrorStack(err))
			return errors.Trace(err)
		case result, ok := <-resultc:
			if !ok {
				resultc = nil
				break
			}
			wg.Add(1)

			// TODO(waigani) this is brittle and expects the result struct to have a
			// DecoratorOptions field. We need to refactor the Flow server to return a
			// tuple of [<decorator>, <result>].
			decorator := reflect.Indirect(reflect.ValueOf(result)).FieldByName("DecoratorOptions").String()
			keep, err := f.ConfirmDecorated(decorator, result)
			if err != nil {
				cancel()
				return errors.Trace(err)
			}
			if keep {

				go func(string, proto.Message) {
					defer wg.Done()
					ctx, err := NewCtx(&f.decoratorApp.App, strings.Split(decorator, " ")[1:])
					if err != nil {
						cancel()
						util.Logger.Fatalf("error getting decorated context: %q", err)
						return
					}
					f.decoratedResultc <- &DecoratedResult{
						Ctx:     ctx,
						Payload: result,
					}
				}(decorator, result)

			}

		case <-timeout:
			cancel()
			return errors.New("timed out waiting for issue")
		}
		if resultc == nil && errc == nil {
			break l
		}
	}
	wg.Done()

	return
}

func (f *flowRunner) CliCtx() (*cli.Context, error) {
	if f.cliCtx == nil {
		ctx, err := NewCtx(&f.cliApp.App, os.Args[1:])
		if err != nil {
			return nil, errors.Trace(err)
		}
		f.cliCtx = ctx
	}
	return f.cliCtx, nil
}

// TODO(waicmdgani) move this to codelingo/sdk/flow
func (f *flowRunner) RunCLI() (chan proto.Message, chan error, func(), error) {
	ctx, err := f.CliCtx()
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}
	if ctx.Bool("debug") {
		err := util.SetDebugLogger()
		if err != nil {
			return nil, nil, nil, errors.Trace(err)
		}
	}

	return f.cliApp.Request(ctx)
}

func (f *flowRunner) ConfirmDecorated(decorator string, payload proto.Message) (bool, error) {
	ctx, err := NewCtx(&f.decoratorApp.App, strings.Split(decorator, " "))
	if err != nil {
		return false, errors.Trace(err)
	}

	return f.decoratorApp.ConfirmDecorated(ctx, payload)
}

func NewCtx(app *cli.App, input []string) (*cli.Context, error) {

	fSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	for _, flag := range app.Flags {
		flag.Apply(fSet)
	}

	if err := fSet.Parse(input); err != nil {
		return nil, errors.Trace(err)
	}
	if err := normalizeFlags(app.Flags, fSet); err != nil {
		return nil, errors.Trace(err)
	}

	ctx := cli.NewContext(nil, fSet, nil)

	return ctx, nil
}

func normalizeFlags(flags []cli.Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.GetName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case *cli.StringSlice:
	default:
		set.Set(name, ff.Value.String())
	}
}
