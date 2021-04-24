package pushover

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMessageClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in -short mode.")
	}

	appToken := os.Getenv("PUSHOVER_TEST_APP_TOKEN")
	if appToken == "" {
		t.Skip("PUSHOVER_TEST_APP_TOKEN is not specified, skipping test.")
	}

	userToken := os.Getenv("PUSHOVER_TEST_USER_TOKEN")
	if userToken == "" {
		t.Skip("PUSHOVER_TEST_USER_TOKEN is not specified, skipping test.")
	}

	ctx := context.Background()

	c, err := NewClient(appToken)
	require.NoError(t, err)

	t.Run("Send", func(t *testing.T) {
		err = c.Send(ctx, userToken, fmt.Sprintf("%s %s", t.Name(), time.Now()))
		require.NoError(t, err)
	})

	t.Run("SendGlance", func(t *testing.T) {
		title := t.Name()
		text := fmt.Sprint(time.Now())
		g := &Glance{
			User:  userToken,
			Title: &title,
			Text:  &text,
		}
		err = c.SendGlance(ctx, g)
		require.NoError(t, err)
	})
}
