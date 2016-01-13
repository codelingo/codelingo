package commands

import (
	"bytes"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/browser"
	"github.com/russross/blackfriday"
)

var DocsCMD = cli.Command{
	Name:   "docs",
	Usage:  "output documentation generated from tenets",
	Action: docs,
}

func docs(c *cli.Context) {

	fmt.Println("Opening tenet documentation in a new browser window ...")

	var mdBuf bytes.Buffer
	if err := writeTenetDoc(c, "", &mdBuf); err != nil {
		oserrf("%v", err)
		return
	}

	// Render markdown to HTML, and sanitise it to protect from
	// malicious tenets.
	htmlUnsafe := blackfriday.MarkdownCommon(mdBuf.Bytes())
	html := bluemonday.UGCPolicy().SanitizeBytes(htmlUnsafe)

	if err := browser.OpenReader(bytes.NewReader(html)); err != nil {
		oserrf("opening docs in browser: %v", err)
		return
	}
}
