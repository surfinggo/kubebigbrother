package channels

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"html/template"
	"io"
	"k8s.io/klog/v2"
	"os"
)

// ChannelPrintConfig is config for ChannelPrint
type ChannelPrintConfig struct {
	Writer          string
	AddedTemplate   string
	DeletedTemplate string
	UpdatedTemplate string
}

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Writer      io.Writer
	IsStdout    bool
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// NewEventProcessContext implements Channel
func (c *ChannelPrint) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  nil,
	}
}

// Handle implements Channel
func (c *ChannelPrint) Handle(ctx *EventProcessContext) error {
	var t *template.Template
	switch ctx.Event.Type {
	case event.TypeAdded:
		t = c.TmplAdded
	case event.TypeDeleted:
		t = c.TmplDeleted
	case event.TypeUpdated:
		t = c.TmplUpdated
	default:
		return errors.Errorf("unknown event type: %s", ctx.Event.Type)
	}

	if c.IsStdout {
		printFunc := func() error {
			buf := &bytes.Buffer{}
			if err := t.Execute(buf, ctx.Event); err != nil {
				// print an extra blank line when error occurs,
				// because print may be interrupted
				// without line feed at the end
				_, _ = c.Writer.Write([]byte("\n"))
				return err
			}

			var styled string
			switch ctx.Event.Type {
			case event.TypeAdded:
				styled = style.Success(buf.String()).String()
			case event.TypeDeleted:
				styled = style.Warning(buf.String()).String()
			default:
				styled = style.Info(buf.String()).String()
			}

			if _, err := c.Writer.Write([]byte(styled)); err != nil {
				return errors.Wrap(err, "write error")
			}
			return nil
		}

		// https://stackoverflow.com/questions/14694088
		return klog.WithLock(printFunc)
	} // end if isStdout

	if err := t.Execute(c.Writer, ctx.Event); err != nil {
		// print an extra blank line when error occurs,
		// because print may be interrupted
		// without line feed at the end
		_, _ = c.Writer.Write([]byte("\n"))
		return errors.Wrap(err, "execute template error")
	}
	return nil
}

const (
	// PrintWriterStdout writes output to stdout
	PrintWriterStdout = "stdout"
)

// NewChannelPrint creates print channel
func NewChannelPrint(config *ChannelPrintConfig) (*ChannelPrint, error) {
	var writer io.Writer
	switch config.Writer {
	case PrintWriterStdout, "":
		writer = os.Stdout
	default:
		return nil, errors.Errorf("unsupported writer: %s", config.Writer)
	}

	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		config.AddedTemplate, config.DeletedTemplate, config.UpdatedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	return &ChannelPrint{
		Writer:      writer,
		IsStdout:    config.Writer == PrintWriterStdout,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
