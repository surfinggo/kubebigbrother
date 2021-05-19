package channels

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"io"
	"os"
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

func NewChannelPrintWithWriter(writer io.Writer) *ChannelPrint {
	return &ChannelPrint{
		Writer: writer,
		// TODO: make WriteFunc configurable
		WriteFunc: func(e *event.Event, w io.Writer) error {
			t := fmt.Sprintf("[%s] %s\n", e.Type, utils.NamespaceKey(e.Obj))
			_, err := w.Write([]byte(t))
			return err
		},
	}

}

const (
	PrintWriterStdout = "stdout"
)

func NewChannelPrint(writerType string) (*ChannelPrint, error) {
	var writer io.Writer
	switch writerType {
	case PrintWriterStdout, "":
		writer = os.Stdout
	default:
		return nil, errors.Errorf("unsupported writer: %s", writerType)
	}
	return NewChannelPrintWithWriter(writer), nil
}
