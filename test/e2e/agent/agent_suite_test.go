package agent_test

import (
	"context"
	"os"
	"testing"

	"github.com/flightctl/flightctl/test/harness/e2e"
	testutil "github.com/flightctl/flightctl/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const TIMEOUT = "5m"
const POLLING = "125ms"
const LONGTIMEOUT = "10m"

// Define a type for messages.
type Message string

const (
	UpdateRenderedVersionSuccess    Message = "Updated to desired renderedVersion: 2"
	ComposeFile                     string  = "podman-compose.yaml"
	ExpectedNumSleepAppV1Containers string  = "3"
	ExpectedNumSleepAppV2Containers string  = "1"
	ZeroContainers                  string  = "0"
)

// String returns the string representation of a message.
func (m Message) String() string {
	return string(m)
}

var (
	suiteCtx     context.Context
	suiteHarness *e2e.Harness
)

func TestAgent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agent E2E Suite")
}

var _ = BeforeSuite(func() {
	suiteCtx = testutil.InitSuiteTracerForGinkgo("Agent E2E Suite")

	// Clean up any existing overlays first if requested
	if os.Getenv("CLEANUP_ALL_SNAPSHOTS") == "true" {
		tempHarness := e2e.NewTestHarness(suiteCtx)
		_ = tempHarness.CleanupAllOverlays()
		tempHarness.Cleanup(false)
	}

	// Create suite-level harness with VM overlay for fast test startup
	suiteHarness = e2e.NewTestHarnessWithOverlay(suiteCtx)
})

var _ = AfterSuite(func() {
	if suiteHarness != nil {
		// Clean up the shared overlay
		if os.Getenv("CLEANUP_ALL_SNAPSHOTS") == "true" {
			suiteHarness.CleanupAllOverlays()
		} else {
			suiteHarness.CleanupOverlays()
		}
		suiteHarness.Cleanup(false)
	}
})
