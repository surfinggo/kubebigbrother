package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	spg "github.com/spongeprojects/client-go/api/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"io"
	"k8s.io/klog/v2"
	"os"
	"text/template"
)

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Writer      io.Writer
	IsStdout    bool
	UseTemplate bool
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
	buf := &bytes.Buffer{}
	if c.UseTemplate {
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
		if err := t.Execute(buf, ctx.Event); err != nil {
			return errors.Wrap(err, "execute template error")
		}
	} else {
		if err := json.NewEncoder(buf).Encode(ctx.Event); err != nil {
			return errors.Wrap(err, "json encode error")
		}
	}

	if c.IsStdout {
		printFunc := func() error {
			styled := style.Fg(ctx.Event.Color(), buf.String()).String()
			if _, err := c.Writer.Write([]byte(styled)); err != nil {
				return errors.Wrap(err, "write error")
			}
			return nil
		}

		// https://stackoverflow.com/questions/14694088
		return klog.WithLock(printFunc)
	} // end if isStdout

	if _, err := io.Copy(c.Writer, buf); err != nil {
		return errors.Wrap(err, "write error")
	}
	return nil
}

const (
	// PrintWriterStdout writes output to stdout
	PrintWriterStdout = "stdout"
)

// NewChannelPrint creates print channel
func NewChannelPrint(config *spg.ChannelPrintConfig) (*ChannelPrint, error) {
	var writer io.Writer
	switch config.Writer {
	case PrintWriterStdout, "":
		klog.V(2).Info("print to stdout")
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
		UseTemplate: config.UseTemplate,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
