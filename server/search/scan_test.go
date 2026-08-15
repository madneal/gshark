package search

import (
	"os"
	"testing"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	global.GVA_LOG = zap.NewNop()
	os.Exit(m.Run())
}

func TestRunProviderWaitsForSlowFnAndKeepsHeartbeat(t *testing.T) {
	stubScanLogPersistence(t)
	origHeartbeat := providerHeartbeat
	providerHeartbeat = time.Millisecond
	heartbeats := 0
	heartbeatScanLog = func(uint, time.Time) error {
		heartbeats++
		return nil
	}
	defer func() {
		providerHeartbeat = origHeartbeat
	}()

	block := make(chan struct{})
	finished := make(chan model.ScanOutcome, 1)

	done := make(chan struct{})
	go func() {
		finished <- runProvider(1, "slow", func() model.ScanOutcome {
			<-block
			return model.ScanSuccess("done")
		})
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("runProvider returned before the provider completed")
	case <-time.After(20 * time.Millisecond):
	}
	close(block)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runProvider did not return after the provider completed")
	}
	if outcome := <-finished; outcome.Status != model.ScanStatusSuccess {
		t.Fatalf("finished status = %q, want success", outcome.Status)
	}
	if heartbeats == 0 {
		t.Fatal("expected heartbeats while the provider was still running")
	}
}

func TestRunProviderReturnsPromptlyOnFastFn(t *testing.T) {
	stubScanLogPersistence(t)
	origHeartbeat := providerHeartbeat
	providerHeartbeat = time.Minute
	defer func() {
		providerHeartbeat = origHeartbeat
	}()

	start := time.Now()
	ran := false
	outcome := runProvider(1, "fast", func() model.ScanOutcome {
		ran = true
		return model.ScanSuccess("done")
	})

	if !ran {
		t.Fatal("expected fn to run")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("runProvider took too long for a fast fn: %s", elapsed)
	}
	if outcome.Status != model.ScanStatusSuccess {
		t.Fatalf("unexpected outcome: %#v", outcome)
	}
}

func TestRunProviderRecordsPanicAsFailure(t *testing.T) {
	stubScanLogPersistence(t)
	origHeartbeat := providerHeartbeat
	providerHeartbeat = time.Minute
	defer func() { providerHeartbeat = origHeartbeat }()

	outcome := runProvider(1, "panic", func() model.ScanOutcome {
		panic("boom")
	})
	if outcome.Status != model.ScanStatusFailed {
		t.Fatalf("outcome status = %q, want %q", outcome.Status, model.ScanStatusFailed)
	}
}

func stubScanLogPersistence(t *testing.T) {
	t.Helper()
	originalStart := startScanLog
	originalHeartbeat := heartbeatScanLog
	originalFinish := finishScanLog
	startScanLog = func(uint, time.Time) error { return nil }
	heartbeatScanLog = func(uint, time.Time) error { return nil }
	finishScanLog = func(uint, model.ScanOutcome, time.Time, time.Time) error { return nil }
	t.Cleanup(func() {
		startScanLog = originalStart
		heartbeatScanLog = originalHeartbeat
		finishScanLog = originalFinish
	})
}
