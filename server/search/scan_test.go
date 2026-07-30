package search

import (
	"os"
	"testing"
	"time"

	"github.com/madneal/gshark/global"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	global.GVA_LOG = zap.NewNop()
	os.Exit(m.Run())
}

func TestRunProviderDoesNotBlockOnSlowFn(t *testing.T) {
	origWatchdog := providerWatchdog
	providerWatchdog = 10 * time.Millisecond
	defer func() { providerWatchdog = origWatchdog }()

	block := make(chan struct{})
	defer close(block)

	done := make(chan struct{})
	go func() {
		runProvider("slow", func() { <-block })
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("runProvider did not return promptly when fn hung past the watchdog")
	}
}

func TestRunProviderReturnsPromptlyOnFastFn(t *testing.T) {
	origWatchdog := providerWatchdog
	providerWatchdog = 1 * time.Minute
	defer func() { providerWatchdog = origWatchdog }()

	start := time.Now()
	ran := false
	runProvider("fast", func() { ran = true })

	if !ran {
		t.Fatal("expected fn to run")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("runProvider took too long for a fast fn: %s", elapsed)
	}
}
