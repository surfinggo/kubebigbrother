package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"html/template"
	"io"
	"k8s.io/klog/v2"
	"os"
)

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Writer    io.Writer
	WriteFunc func(*event.Event, io.Writer) error
}

// Handle implements Channel
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

func NewChannelPrintWithWriter(writer io.Writer, isStdout bool,
	addedTmpl, deletedTmpl, updatedTmpl string) (*ChannelPrint, error) {
	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		addedTmpl, deletedTmpl, updatedTmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	writeFunc := func(e *event.Event, w io.Writer) error {
		var t *template.Template
		switch e.Type {
		case event.TypeAdded:
			t = tmplAdded
		case event.TypeDeleted:
			t = tmplDeleted
		case event.TypeUpdated:
			t = tmplUpdated
		default:
			panic(fmt.Sprintf("unknown event type: %s", e.Type))
		}

		if isStdout {
			printFunc := func() error {
				buf := &bytes.Buffer{}
				if err := t.Execute(buf, e); err != nil {
					// print an extra blank line when error occurs,
					// because print may be interrupted
					// without line feed at the end
					_, _ = w.Write([]byte("\n"))
					return err
				}
				var styled string
				switch e.Type {
				case event.TypeAdded:
					styled = style.Success(buf.String()).String()
				case event.TypeDeleted:
					styled = style.Warning(buf.String()).String()
				default:
					styled = style.Info(buf.String()).String()
				}

				if _, err := w.Write([]byte(styled)); err != nil {
					return errors.Wrap(err, "write to writer error")
				}

				if err := t.Execute(w, e); err != nil {
					// print an extra blank line when error occurs,
					// because print may be interrupted
					// without line feed at the end
					_, _ = w.Write([]byte("\n"))
					return errors.Wrap(err, "execute template error")
				}
				return nil
			}

			// https://stackoverflow.com/questions/14694088
			return klog.WithLock(printFunc)
		} // end if isStdout

		if err := t.Execute(w, e); err != nil {
			// print an extra blank line when error occurs,
			// because print may be interrupted
			// without line feed at the end
			_, _ = w.Write([]byte("\n"))
			return errors.Wrap(err, "execute template error")
		}
		return nil
	} // end writeFunc

	return &ChannelPrint{
		Writer:    writer,
		WriteFunc: writeFunc,
	}, nil
}

const (
	PrintWriterStdout = "stdout"
)

func NewChannelPrint(writerType,
	addedTmpl, deletedTmpl, updatedTmpl string) (*ChannelPrint, error) {
	var writer io.Writer
	switch writerType {
	case PrintWriterStdout, "":
		writer = os.Stdout
	default:
		return nil, errors.Errorf("unsupported writer: %s", writerType)
	}
	return NewChannelPrintWithWriter(writer, writerType == PrintWriterStdout, addedTmpl, deletedTmpl, updatedTmpl)
}
