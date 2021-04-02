package ic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"github.com/shurcooL/go/reflectsource"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	prefixFunction = func() string {
		return "ic| "
	}
	defaultPrefixFunction = prefixFunction

	includeContext = false

	outputFunction = func (x string) {
		_, err := fmt.Fprintln(os.Stderr, x)
		if err != nil {
			panic(err)
		}
	}
	defaultOutputFunction = outputFunction
	enableOutput = true

	enableSyntaxHighlighting = true
)

// IC takes an arbitrary number of arbitrary type arguments, formats and optionally syntax highlights them then outputs
// them, by default, on os.Stderr
func IC(o ...interface{}) {
	if enableOutput {
		as, _ := toArgSlice(o, reflectsource.GetParentArgExprAllAsString())
		outputFunction(formatToString(as, enableSyntaxHighlighting))
	}
}

// Format is the same as IC, except it returns the generated string without syntax highlighting instead of outputting it
func Format(o ...interface{}) string {
	as, _ := toArgSlice(o, reflectsource.GetParentArgExprAllAsString())
	return formatToString(as, false)
}

// Enable enables output from IC
func Enable() {
	enableOutput = true
}

// Disable disables output from IC
func Disable() {
	enableOutput = false
}

// ConfigureSetIncludeContext takes a boolean to signify if every call of Format or IC should include a context string
// (which include the source file name, line that contains the IC or Format function call and the package name)
func ConfigureSetIncludeContext(x bool) {
	includeContext = x
}

// ConfigureResetPrefix resets the prefix for outputs from IC and Format to default ("ic| "). This overrides anything
// set in ConfigureSetPrefix or ConfigureSetPrefixFunction
func ConfigureResetPrefix() {
	prefixFunction = defaultPrefixFunction
}

// ConfigureSetPrefix sets the prefix for output from IC and Format. Can be overridden by subsequent calls to
// ConfigureResetPrefix or ConfigureSetPrefixFunction
func ConfigureSetPrefix(newPrefix string) {
	prefixFunction = func() string {
		return newPrefix
	}
}

// ConfigureSetPrefixFunction sets the prefix for the output from IC and Format to the result of the supplied function.
// Can be overridden by subsequent calls to ConfigureResetPrefix or ConfigureSetPrefix
func ConfigureSetPrefixFunction(pf func() string) {
	prefixFunction = pf
}

// ConfigureResetOutput resets the output of IC to os.Stderr. Can be overridden by subsequent calls to
// ConfigureSetOutput or ConfigureSetOutputFunction
func ConfigureResetOutput() {
	outputFunction = defaultOutputFunction
}

// ConfigureSetOutput sets the output of IC to the provided io.Writer. Can be overridden by subsequent calls to
// ConfigureResetOutput or ConfigureSetOutputFunction. A panic will occur if the supplied writer cannot be written to.
func ConfigureSetOutput(wr io.Writer) {
	outputFunction = func(x string) {
		_, err := fmt.Fprintln(wr, x)
		if err != nil {
			panic(err)
		}
	}
}

// ConfigureSetOutputFunction sets the output of IC to an arbitrary function. Can be overridden by subsequent calls to
// ConfigureResetOutput or ConfigureSetOutput
func ConfigureSetOutputFunction(x func (string)) {
	outputFunction = x
}

// ConfigureEnableSyntaxHighlighting enables syntax highlighting on output from IC
func ConfigureEnableSyntaxHighlighting() {
	enableSyntaxHighlighting = true
}

// ConfigureDisableSyntaxHighlighting disables syntax highlighting on output from IC
func ConfigureDisableSyntaxHighlighting() {
	enableSyntaxHighlighting = false
}

type argument struct {
	Source string
	Value interface{}
}

func toArgSlice(o []interface{}, valSources []string) ([]*argument, error) {
	if len(o) != len(valSources) {
		// theoretically this should never occur
		// if it does, you'll just get an empty slice which will trigger a message like "ic| filename.go:123 in blah"
		return nil, errors.New("toArgSlice: length of both arguments must be identical")
	}

	var avs []*argument
	for i, v := range o {
		avs = append(avs, &argument{
			Source: valSources[i],
			Value:  v,
		})
	}

	return avs, nil
}

func highlight(s string, enable bool) string {
	if !enable {
		return s
	}
	b := new(bytes.Buffer)
	_ = quick.Highlight(b, s, "go", "terminal256", "monokai")
	x := b.String()
	return strings.ReplaceAll(x, "\n", "") // Chroma puts some really weird newlines in the output
	// since reflectsource doesn't actually work properly with arguments that span multiple lines, I think we'll be okay
	// doing this
}

func makeContext() string {
	pc, file, no, ok := runtime.Caller(3)
	if ok {

		fname := runtime.FuncForPC(pc).Name()
		var packageName string
		{
			firstDot := strings.Index(fname, ".")
			if firstDot == -1 {
				packageName = "<unknown>"
			} else {
				packageName = fname
			}
		}
		_, filename := filepath.Split(file)

		return fmt.Sprintf("%s:%d in %s", filename, no, packageName)
	}
	return "unable to determine caller"
}

func formatToString(avs []*argument, doSyntaxHighlighting bool) string {
	var outputParts []string

	if len(avs) == 0 {
		outputParts = append(outputParts, makeContext())
	} else {

		for i, av := range avs {

			var x string

			switch v := av.Value.(type) {
			case string:
				quoted := strconv.Quote(v)
				if quoted == av.Source {
					x = highlight(av.Source, doSyntaxHighlighting)
				} else {
					x = fmt.Sprintf("%s: %s", highlight(av.Source, doSyntaxHighlighting), highlight(quoted, doSyntaxHighlighting))
				}
			default:
				if av.Source == fmt.Sprint(av.Value) {
					x = highlight(fmt.Sprint(av.Value), doSyntaxHighlighting)
				} else {
					x = fmt.Sprintf("%s: %s", highlight(av.Source, doSyntaxHighlighting), highlight(fmt.Sprintf("%+v", v), doSyntaxHighlighting))
				}
			}

			if includeContext {
				if len(avs) > 1 {
					if i == 0 {
						outputParts = append(outputParts, makeContext())
					}
				} else {
					x = fmt.Sprintf("%s - %s", makeContext(), x)
				}
			}

			outputParts = append(outputParts, x)

		}

	}

	stringPrefix := prefixFunction()
	prefixLen := len(stringPrefix)
	for i, item := range outputParts {
		if i == 0 {
			item = stringPrefix + item
		} else {
			item = strings.Repeat(" ", prefixLen) + item
		}
		outputParts[i] = item
	}

	return strings.Join(outputParts, "\n")
}