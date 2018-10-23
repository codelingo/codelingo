package flow

// TODO(waigani) currently hardcoded to rewrite. Generalise this so it can be
// used for all Flow endpoints. It may be possible to use the command config
// structs to do this.

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/codelingo/lingo/app/util"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

type Confirmer struct {
	cancel      context.CancelFunc
	ctx         *cli.Context
	msgc        chan proto.Message
	errorc      chan error
	itemFactory func(msg proto.Message) *ConfirmerItem
}

func (c *Confirmer) Confirm() (confirmed []proto.Message, err error) {

	defer func() {
		if err != nil {
			c.cancel()
		}
	}()
	defer util.Logger.Sync()

	// spnr := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	// spnr.Start()
	// defer spnr.Stop()

	keepAll := c.ctx.IsSet("keep-all")
	// If user is manually confirming reviews, set a long timeout.
	timeout := time.After(time.Hour * 1)
	if keepAll {
		timeout = time.After(time.Minute * 11)
	}

l:
	for {
		select {
		case err, ok := <-c.errorc:
			if !ok {
				c.errorc = nil
				break
			}

			util.Logger.Debugf("error: %s", errors.ErrorStack(err))
			return nil, errors.Trace(err)
		case msg, ok := <-c.msgc:
			// if !keepAll {
			// 	spnr.Stop()
			// }
			if !ok {
				c.msgc = nil
				break
			}
			util.Logger.Debugf("received msg %v", msg)

			item := c.itemFactory(msg)

			pass, err := item.Confirm(c.ctx)
			if err != nil {
				return nil, errors.Trace(err)
			}
			if pass {
				confirmed = append(confirmed, msg)
			}

			// if !keepAll {
			// 	spnr.Restart()
			// }
		case <-timeout:
			c.cancel()
			return nil, errors.New("timed out waiting for response")
		}
		if c.msgc == nil && c.errorc == nil {
			break l
		}
	}

	// Stop spinner if it hasn't been stopped already
	// if keepAll {
	// 	spnr.Stop()
	// }
	return
}

type ConfirmerItem struct {
	attempt int
	Preview func() string

	// Options is a map of option keys to confirm functions. Each confirm function returns: <keep>bool, <retry>bool, <err>error
	Options map[string]func() (bool, bool, error)

	// OptionKeyMap maps an option key to aliases e.g. "[k]eep" => "k", "keep", "K"
	OptionKeyMap map[string][]string

	// FlagOptions is a map of flags to options show if that flag is present.
	// "_all" is a builtin flag for those options that don't require a flag.
	FlagOptions map[string][]string
}

// confirm returns true if the msg should be kept or false if it should be
// dropped.
func (item *ConfirmerItem) Confirm(ctx *cli.Context) (bool, error) {

	if item.attempt == 0 {
		fmt.Println(item.Preview())
	}
	item.attempt++

	var option string
	fmt.Print("\n")

	for flag, opts := range item.FlagOptions {
		if flag == "_all" || ctx.IsSet(flag) {
			for _, opt := range opts {
				fmt.Print(opt, " ")
			}
		}
	}
	fmt.Print(": ")
	fmt.Scanln(&option)
	action, err := item.action(option)
	if err != nil {
		return false, errors.Trace(err)
	}

	pass, retry, err := action()
	if err != nil {
		return false, errors.Trace(err)
	}
	if retry {
		return item.Confirm(ctx)
	}
	return pass, err
}

func (c *ConfirmerItem) action(option string) (func() (bool, bool, error), error) {
	for optKey, aliases := range c.OptionKeyMap {
		for _, alias := range aliases {
			if alias == option {
				return c.Options[optKey], nil
			}
		}
	}

	return nil, errors.Errorf("no action found for option %q", option)
}

// confirm actions

func OpenFileConfirmAction(filename string, line int64) (bool, bool, error) {
	var editor string
	var app string
	defaultEditor := "vi" // TODO(waigani) use EDITOR or VISUAL env vars
	// https://github.com/codelingo/demo/msgs/3
	if editor != "" {
		defaultEditor = editor
	}
	fmt.Printf("application (%s):", defaultEditor)
	fmt.Scanln(&app)
	if app == "" {
		app = defaultEditor
	}
	// c := msg.Position.Start.Column // TODO(waigani) use column
	// https://github.com/codelingo/demo/msgs/4

	cmd, err := util.OpenFileCmd(app, filename, line)
	if err != nil {
		return false, true, err
	}

	if err = cmd.Start(); err != nil {
		log.Println(err)
	}
	if err = cmd.Wait(); err != nil {
		log.Println(err)
	}

	editor = app
	return false, true, nil
}
