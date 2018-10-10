package render

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/samfoo/ansi"
)

func Terminal() (blackfriday.Renderer, int) {
	extensions := 0 |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_DEFINITION_LISTS
	return &Console{}, extensions
}

// Index corresponds to the heading level (e.g. h1, h2, h3...)
var headerStyles = [...]string{
	ansi.ColorCode("yellow+bhu"),
	ansi.ColorCode("yellow+bh"),
	ansi.ColorCode("yellow"),
	ansi.ColorCode("yellow"),
}

var emphasisStyles = [...]string{
	ansi.ColorCode("cyan+bh"),
	ansi.ColorCode("cyan+bhu"),
	ansi.ColorCode("cyan+bhi"),
}

var linkStyle = ansi.ColorCode("015+u")

const (
	UNORDERED = 1 << iota
	ORDERED
)

type list struct {
	kind  int
	index int
}

type Console struct {
	lists []*list
}

func (options *Console) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	s := string(text)
	reg, _ := regexp.Compile("\\n")

	out.WriteString("\n    ")
	out.WriteString(ansi.ColorCode("015"))
	out.WriteString(reg.ReplaceAllString(s, "\n    "))
	out.WriteString(ansi.ColorCode("reset"))
	out.WriteString("\n")
}

func (options *Console) BlockQuote(out *bytes.Buffer, text []byte) {
	s := strings.TrimSpace(string(text))
	reg, _ := regexp.Compile("\\n")

	out.WriteString("\n  | ")
	out.WriteString(reg.ReplaceAllString(s, "\n  | "))
	out.WriteString("\n\n")
}

func (options *Console) BlockHtml(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *Console) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	out.WriteString("\n")
	out.WriteString(headerStyles[level-1])

	marker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}

	out.WriteString(ansi.ColorCode("reset"))
	out.WriteString("\n\n")
}

func (options *Console) HRule(out *bytes.Buffer) {
	out.WriteString("\n\u2015\u2015\u2015\u2015\u2015\n\n")
}

func (options *Console) List(out *bytes.Buffer, text func() bool, flags int) {
	out.WriteString("\n")

	kind := UNORDERED
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		kind = ORDERED
	}

	options.lists = append(options.lists, &list{kind, 1})
	text()
	options.lists = options.lists[:len(options.lists)-1]
	out.WriteString("\n")
}

func (options *Console) ListItem(out *bytes.Buffer, text []byte, flags int) {
	current := options.lists[len(options.lists)-1]

	for i := 0; i < len(options.lists); i++ {
		out.WriteString("  ")
	}

	if current.kind == ORDERED {
		out.WriteString(fmt.Sprintf("%d. ", current.index))
		current.index += 1
	} else {
		out.WriteString(ansi.ColorCode("red+bh"))
		out.WriteString("* ")
		out.WriteString(ansi.ColorCode("reset"))
	}

	out.Write(text)
	out.WriteString("\n\n")
}

func (options *Console) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
		return
	}

	out.WriteString("\n\n")
}

func (options *Console) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (options *Console) TableRow(out *bytes.Buffer, text []byte)                               {}
func (options *Console) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             {}
func (options *Console) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (options *Console) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (options *Console) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}

func (options *Console) TitleBlock(out *bytes.Buffer, text []byte) {
	out.WriteString("\n")
	out.WriteString(headerStyles[0])
	out.Write(text)
	out.WriteString(ansi.ColorCode("reset"))
	out.WriteString("\n\n")
}

func (options *Console) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString(linkStyle)
	out.Write(link)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString(ansi.ColorCode("015+b"))
	out.Write(text)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString(emphasisStyles[1])
	out.Write(text)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString(emphasisStyles[0])
	out.Write(text)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	out.WriteString(" [ image ] ")
}

func (options *Console) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n")
}

func (options *Console) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.Write(content)
	out.WriteString(" (")
	out.WriteString(linkStyle)
	out.Write(link)
	out.WriteString(ansi.ColorCode("reset"))
	out.WriteString(")")
}

func (options *Console) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	out.WriteString(ansi.ColorCode("magenta"))
	out.Write(tag)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString(emphasisStyles[2])
	out.Write(text)
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString(ansi.ColorCode("008+s"))
	out.WriteString("\u2015")
	out.Write(text)
	out.WriteString("\u2015")
	out.WriteString(ansi.ColorCode("reset"))
}

func (options *Console) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
}

func (options *Console) Entity(out *bytes.Buffer, entity []byte) {
	out.WriteString(html.UnescapeString(string(entity)))
}

func (options *Console) NormalText(out *bytes.Buffer, text []byte) {
	s := string(text)
	reg, _ := regexp.Compile("\\s+")

	out.WriteString(reg.ReplaceAllString(s, " "))
}

func (options *Console) DocumentHeader(out *bytes.Buffer) {
}

func (options *Console) DocumentFooter(out *bytes.Buffer) {
}

func (options *Console) GetFlags() int {
	return 0
}
