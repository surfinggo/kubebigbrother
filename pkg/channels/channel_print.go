package channels

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"html/template"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"strings"
)

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Writer    io.Writer
	WriteFunc func(*event.Event, io.Writer) error
}

func (c *ChannelPrint) Handle(e *event.Event) error {
	if c.WriteFunc != nil {
		return c.WriteFunc(e, c.Writer)
	}
	err := json.NewEncoder(c.Writer).Encode(e)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	return nil
}

func NewChannelPrintWithWriter(writer io.Writer, tmpl string) (*ChannelPrint, error) {
	if tmpl == "" {
		tmpl = "[{{.Obj.GroupVersionKind}} {{.Type}}] {{.Obj.GetNamespace}}/{{.Obj.GetName}}\n"
		// example of using field:
		//tmpl = "[{{.Obj.GroupVersionKind}} {{.Type}}] {{.Obj.GetNamespace}}/{{.Obj.GetName}} {{field .Obj \"kind\"}}\n"
	}
	funcMap := template.FuncMap{
		"field": func(s *unstructured.Unstructured, path ...string) string {
			// methods can be used in template:
			// s.GetName()
			// s.GetNamespace()
			str, exist, err := unstructured.NestedString(s.Object, path...)
			if err != nil {
				return fmt.Sprintf("[Error reading field .%s: %s]", strings.Join(path, "."), err)
			}
			if !exist {
				return fmt.Sprintf("[Field .%s not exist]", strings.Join(path, "."))
			}
			return str
		},
	}
	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	return &ChannelPrint{
		Writer: writer,
		WriteFunc: func(e *event.Event, w io.Writer) error {
			err := t.Execute(w, e)
			if err != nil {
				// print an extra blank line when error occurs,
				// because print may be interrupted
				// without line feed at the end
				_, _ = w.Write([]byte("\n"))
			}
			return err
		},
	}, nil
}

const (
	PrintWriterStdout = "stdout"
)

func NewChannelPrint(writerType, tmpl string) (*ChannelPrint, error) {
	var writer io.Writer
	switch writerType {
	case PrintWriterStdout, "":
		writer = os.Stdout
	default:
		return nil, errors.Errorf("unsupported writer: %s", writerType)
	}
	return NewChannelPrintWithWriter(writer, tmpl)
}
