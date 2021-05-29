package event

import (
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEvent_Color(t *testing.T) {
	assertions := require.New(t)

	e := Event{
		Type: TypeAdded,
	}

	assertions.Equal(style.Green, e.Color())
}
