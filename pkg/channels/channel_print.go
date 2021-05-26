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

func NewChannelPrintWithWriter(writer io.Writer,
	addedTmpl, updatedTmpl, deletedTmpl string) (*ChannelPrint, error) {
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
	if addedTmpl == "" {
		addedTmpl = "[{{.Obj.GroupVersionKind}} is created] {{.Obj.GetNamespace}}/{{.Obj.GetName}}\n"
		// example of using field:
		//tmpl = "[{{.Obj.GroupVersionKind}} is created] " +
		// "{{.Obj.GetNamespace}}/{{.Obj.GetName}} {{field .Obj \"kind\"}}\n"
	}
	if updatedTmpl == "" {
		updatedTmpl = "[{{.Obj.GroupVersionKind}} is updated] {{.Obj.GetNamespace}}/{{.Obj.GetName}}\n"
	}
	if deletedTmpl == "" {
		deletedTmpl = "[{{.Obj.GroupVersionKind}} is deleted] {{.Obj.GetNamespace}}/{{.Obj.GetName}}\n"
	}
	tmplAdded, err := template.New("").Funcs(funcMap).Parse(addedTmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parse added template error")
	}
	tmplUpdated, err := template.New("").Funcs(funcMap).Parse(updatedTmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parse updated template error")
	}
	tmplDeleted, err := template.New("").Funcs(funcMap).Parse(deletedTmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parse deleted template error")
	}

	return &ChannelPrint{
		Writer: writer,
		WriteFunc: func(e *event.Event, w io.Writer) error {
			var t *template.Template
			switch e.Type {
			case event.TypeAdded:
				t = tmplAdded
			case event.TypeUpdated:
				t = tmplUpdated
			case event.TypeDeleted:
				t = tmplDeleted
			default:
				panic(fmt.Sprintf("unknown event type: %s", e.Type))
			}
			if err := t.Execute(w, e); err != nil {
				// print an extra blank line when error occurs,
				// because print may be interrupted
				// without line feed at the end
				_, _ = w.Write([]byte("\n"))
				return err
			}
			return nil
		},
	}, nil
}

const (
	PrintWriterStdout = "stdout"
)

func NewChannelPrint(writerType,
	addedTmpl, updatedTmpl, deletedTmpl string) (*ChannelPrint, error) {
	var writer io.Writer
	switch writerType {
	case PrintWriterStdout, "":
		writer = os.Stdout
	default:
		return nil, errors.Errorf("unsupported writer: %s", writerType)
	}
	return NewChannelPrintWithWriter(writer, addedTmpl, updatedTmpl, deletedTmpl)
}
