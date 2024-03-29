package riemann_listener

import (
	"log"
	"testing"
	"time"

	riemanngo "github.com/riemann/riemann-go-client"
	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/testutil"
)

func TestSocketListener_tcp(t *testing.T) {
	log.Println("Entering")

	sl := newRiemannSocketListener()
	sl.Log = testutil.Logger{}
	sl.ServiceAddress = "tcp://127.0.0.1:5555"
	sl.ReadBufferSize = config.Size(1024)

	acc := &testutil.Accumulator{}
	err := sl.Start(acc)
	require.NoError(t, err)
	defer sl.Stop()

	// Check for stats of specific service
	c := riemanngo.NewTCPClient("127.0.0.1:5555", 5*time.Second)
	require.NoError(t, c.Connect())
	require.NoError(t, err)
	defer c.Close()

	result, err := riemanngo.SendEvent(c, &riemanngo.Event{
		Service: "hello",
	})
	require.NoError(t, err)
	require.True(t, result.GetOk())

	// Check for stats no service specified
	c = riemanngo.NewTCPClient("127.0.0.1:5555", 5*time.Second)
	require.NoError(t, c.Connect())
	defer c.Close()

	result, err = riemanngo.SendEvent(c, &riemanngo.Event{})
	require.False(t, result.GetOk())
	require.Equal(t, "No Service Name", result.GetError())
	require.NoError(t, err)
}
