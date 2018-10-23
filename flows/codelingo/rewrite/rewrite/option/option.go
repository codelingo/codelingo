package option

import "github.com/urfave/cli"

type option uint32

const (
	opt_replace option = 1 << iota
	opt_prepend
	opt_append
	opt_line
	opt_byte
	opt_startOffset
	opt_endOffset
	opt_startToEndOffset
)

func (o option) String() string {
	switch o {
	case opt_replace:
		return "replace"
	case opt_prepend:
		return "prepend"
	case opt_append:
		return "append"
	case opt_line:
		return "line"
	case opt_byte:
		return "byte"
	case opt_startOffset:
		return "start-offset"
	case opt_endOffset:
		return "end-offset"
	case opt_startToEndOffset:
		return "start-to-end-offset"
	}
	return "unknown"
}

func strToOpt(str string) option {
	switch str {
	case "replace":
		return opt_replace
	case "prepend":
		return opt_prepend
	case "append":
		return opt_append
	case "line":
		return opt_line
	case "byte":
		return opt_byte
	case "start-offset":
		return opt_startOffset
	case "end-offset":
		return opt_endOffset
	case "start-to-end-offset":
		return opt_startToEndOffset
	}
	return 0
}

// TODO(waigani) make this a default pattern for all flow decorators

// New returns a new option object built from a decorator string with the
// following format:  <flow>.[<options>].<config_key>
func New(ctx *cli.Context) (option, error) {
	opts := opt_replace | opt_startToEndOffset | opt_byte

	for _, flag := range []string{"replace", "prepend", "append", "line", "byte", "start-offset", "end-offset", "start-to-end-offset"} {
		if ctx.Bool(flag) {
			opts.addOption(strToOpt(flag))
		}
	}

	return opts, nil
}

func (o option) IsReplace() bool {
	return o.mode() == opt_replace
}

func (o option) IsPrepend() bool {
	return o.mode() == opt_prepend
}

func (o option) IsAppend() bool {
	return o.mode() == opt_append
}

func (o option) IsStartToEndOffset() bool {
	return o.offset() == opt_startToEndOffset
}

func (o option) IsStartOffset() bool {
	return o.offset() == opt_startOffset
}

func (o option) IsEndOffset() bool {
	return o.offset() == opt_endOffset
}

func (o option) IsLine() bool {
	return o.position() == opt_line
}

func (o option) IsByte() bool {
	return o.position() == opt_byte
}

// Mode returns the rewrite mode, defaulting to replace if append or prepend
// are not set.
func (o option) mode() option {
	switch {
	case o.hasOption(opt_append):
		return opt_append
	case o.hasOption(opt_prepend):
		return opt_prepend
	}
	return opt_replace
}

// Offset returns the offset(s) to rewrite, defaulting to start-to-end-offset if start-offset or end-offset
// are not set.
func (o option) offset() option {
	switch {
	case o.hasOption(opt_startOffset):
		return opt_startOffset
	case o.hasOption(opt_endOffset):
		return opt_endOffset
	}
	return opt_startToEndOffset
}

// Position returns the position to rewrite, defaulting to byte if line is not set.
func (o option) position() option {
	if o.hasOption(opt_line) {
		return opt_line
	}
	return opt_byte
}

func (o option) hasOption(opt option) bool {
	return o&opt != 0
}

func (o *option) addOption(opt option) {
	*o |= opt
}

func (o *option) clearOption(opt option) {
	*o &= ^opt
}

func (o *option) toggleOption(opt option) {
	*o ^= opt
}
