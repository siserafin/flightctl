package microshift

import (
	"context"
	"os"
	"testing"

	"github.com/flightctl/flightctl/test/harness/e2e"
	testutil "github.com/flightctl/flightctl/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	suiteCtx     context.Context
	suiteHarness *e2e.Harness
)

func TestMicroshift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "microshift ACM enrollment E2E Suite")
}

var _ = BeforeSuite(func() {
	suiteCtx = testutil.InitSuiteTracerForGinkgo("microshift ACM enrollment E2E Suite")

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
