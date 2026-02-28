package rollout_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/flightctl/flightctl/test/harness/e2e"
	testutil "github.com/flightctl/flightctl/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const TIMEOUT = "5m"
const POLLING = "125ms"
const FASTPOLLING = "100ms" // Fast polling for catching quick batch transitions
const POLLINGINTERVAL = "10s"
const MEDIUMTIMEOUT = "10m"
const LONGTIMEOUT = "15m"
const DEVICEWAITTIME = "30s"
const DEFAULTUPDATETIMEOUT = "90s"

func TestRollout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rollout Suite")
}

var _ = BeforeSuite(func() {
	// Clean up ALL leftover e2e VMs and temp directories from previous test runs
	// This ensures a completely clean slate before rollout tests start
	GinkgoWriter.Printf("🔄 [BeforeSuite] Rollout: Running global e2e cleanup from previous runs\n")
	cleanupFreshVMs()

	// Setup harness without VM for rollout tests
	// Rollout tests only need API access, device VMs are created separately with worker IDs 1000+
	GinkgoWriter.Printf("🔄 [BeforeSuite] Rollout: Starting harness setup without VM\n")
	_, _, err := e2e.SetupWorkerHarnessWithoutVM()
	if err != nil {
		GinkgoWriter.Printf("❌ [BeforeSuite] Rollout: Failed to setup harness: %v\n", err)
		Expect(err).ToNot(HaveOccurred())
	}
	GinkgoWriter.Printf("✅ [BeforeSuite] Rollout: Harness setup completed successfully\n")
})

var _ = BeforeEach(func() {
	// Get the harness and context directly - no package-level variables
	workerID := GinkgoParallelProcess()
	harness := e2e.GetWorkerHarness()
	suiteCtx := e2e.GetWorkerContext()

	GinkgoWriter.Printf("🔄 [BeforeEach] Worker %d: Setting up test (no VM needed for main harness)\n", workerID)

	// Create test-specific context for proper tracing
	testCtx := testutil.StartSpecTracerForGinkgo(suiteCtx)

	// Set the test context in the harness
	harness.SetTestContext(testCtx)

	// No VM setup needed - rollout tests use device VMs (worker IDs 1000+) only
	GinkgoWriter.Printf("✅ [BeforeEach] Worker %d: Test setup completed\n", workerID)
})

var _ = AfterEach(func() {
	workerID := GinkgoParallelProcess()
	GinkgoWriter.Printf("🔄 [AfterEach] Worker %d: Cleaning up test resources\n", workerID)

	// Get the harness and context directly - no shared variables needed
	harness := e2e.GetWorkerHarness()
	suiteCtx := e2e.GetWorkerContext()

	// Clean up test resources BEFORE switching back to suite context
	// This ensures we use the correct test ID for resource cleanup
	err := harness.CleanUpAllTestResources()
	Expect(err).ToNot(HaveOccurred())

	// Now restore suite context for any remaining cleanup operations
	harness.SetTestContext(suiteCtx)

	GinkgoWriter.Printf("✅ [AfterEach] Worker %d: Test cleanup completed\n", workerID)
})

var _ = AfterSuite(func() {
	workerID := GinkgoParallelProcess()
	GinkgoWriter.Printf("🔄 [AfterSuite] Worker %d: Cleaning up rollout test suite\n", workerID)

	// Clean up ALL e2e VMs and temp directories (same as BeforeSuite)
	cleanupFreshVMs()

	// Also remove VMs from the pool
	for i := 0; i < 10; i++ { // Support up to 10 device VMs per worker
		deviceWorkerID := 1000 + i
		err := e2e.RemoveVMFromPool(deviceWorkerID)
		if err != nil {
			GinkgoWriter.Printf("⚠️  [AfterSuite] Warning: Failed to remove VM for worker %d: %v\n", deviceWorkerID, err)
		}
	}

	GinkgoWriter.Printf("✅ [AfterSuite] Worker %d: Rollout suite cleanup completed\n", workerID)
})

// cleanupFreshVMs cleans up ALL e2e VMs and temporary directories
// This is the same cleanup logic used in test/scripts/e2e_cleanup.sh
// It ensures a completely clean slate before rollout tests start
func cleanupFreshVMs() {
	GinkgoWriter.Printf("🔄 [Cleanup] Starting global E2E test cleanup before rollout tests...\n")
	GinkgoWriter.Printf("🔄 [Cleanup] Finding flightctl e2e VMs...\n")

	// Find all flightctl e2e VMs using virsh
	cmd := exec.Command("virsh", "list", "--all", "--name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		GinkgoWriter.Printf("⚠️  [Cleanup] Failed to list VMs: %v (virsh may not be available)\n", err)
		return
	}

	// Filter for ALL flightctl e2e VMs (includes both pool VMs and imagebuild test VMs)
	// Matches: flightctl-e2e-* and imagebuild-test-*
	vmNames := strings.Split(string(output), "\n")
	flightctlVMs := []string{}
	for _, vmName := range vmNames {
		vmName = strings.TrimSpace(vmName)
		if vmName != "" && (strings.HasPrefix(vmName, "flightctl-e2e-") || strings.HasPrefix(vmName, "imagebuild-test-")) {
			flightctlVMs = append(flightctlVMs, vmName)
		}
	}

	if len(flightctlVMs) == 0 {
		GinkgoWriter.Printf("✅ [Cleanup] No flightctl e2e VMs found to clean up\n")
	} else {
		GinkgoWriter.Printf("🔍 [Cleanup] Found %d flightctl e2e VM(s): %v\n", len(flightctlVMs), flightctlVMs)
	}

	// Clean up each VM
	for _, vmName := range flightctlVMs {
		GinkgoWriter.Printf("🔄 [Cleanup] Cleaning up VM: %s\n", vmName)

		// 1. Delete pristine snapshot (ignore errors if it doesn't exist)
		GinkgoWriter.Printf("   🔄 Deleting pristine snapshot for %s\n", vmName)
		cmd := exec.Command("virsh", "snapshot-delete", vmName, "pristine", "--metadata")
		if err := cmd.Run(); err != nil {
			GinkgoWriter.Printf("   ⚠️  Failed to delete pristine snapshot for %s (may not exist)\n", vmName)
		}

		// 2. Destroy the VM if it's running
		cmd = exec.Command("virsh", "destroy", vmName)
		if err := cmd.Run(); err != nil {
			GinkgoWriter.Printf("   ⚠️  Failed to destroy %s (may not be running)\n", vmName)
		}

		// 3. Undefine the domain (try multiple approaches)
		var undefined bool
		
		// Try simple undefine first
		cmd = exec.Command("virsh", "undefine", vmName)
		if err := cmd.Run(); err == nil {
			GinkgoWriter.Printf("   ✅ Successfully cleaned up VM: %s\n", vmName)
			undefined = true
		}

		// Try with NVRAM
		if !undefined {
			cmd = exec.Command("virsh", "undefine", vmName, "--nvram")
			if err := cmd.Run(); err == nil {
				GinkgoWriter.Printf("   ✅ Successfully cleaned up VM: %s (with NVRAM)\n", vmName)
				undefined = true
			}
		}

		// Try with storage and NVRAM
		if !undefined {
			cmd = exec.Command("virsh", "undefine", vmName, "--remove-all-storage", "--nvram")
			if err := cmd.Run(); err == nil {
				GinkgoWriter.Printf("   ✅ Successfully cleaned up VM: %s (with storage and NVRAM)\n", vmName)
				undefined = true
			}
		}

		if !undefined {
			GinkgoWriter.Printf("   ❌ Failed to undefine %s with all approaches\n", vmName)
		}
	}

	// Clean up temporary directories in /tmp
	GinkgoWriter.Printf("🔄 [Cleanup] Cleaning up temporary directories...\n")
	cmd = exec.Command("find", "/tmp", "-maxdepth", "1", "-name", "flightctl-e2e-worker-*", "-type", "d")
	output, err = cmd.CombinedOutput()
	if err != nil {
		GinkgoWriter.Printf("⚠️  [Cleanup] Failed to search for temporary directories: %v\n", err)
		return
	}

	tmpDirs := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tmpDirs) == 1 && tmpDirs[0] == "" {
		GinkgoWriter.Printf("✅ [Cleanup] No temporary directories found\n")
	} else {
		GinkgoWriter.Printf("🔍 [Cleanup] Found %d temporary director(ies)\n", len(tmpDirs))
		for _, dir := range tmpDirs {
			if dir != "" {
				GinkgoWriter.Printf("   🗑️  Removing: %s\n", dir)
				if err := os.RemoveAll(dir); err != nil {
					GinkgoWriter.Printf("   ⚠️  Failed to remove %s: %v\n", dir, err)
				}
			}
		}
		GinkgoWriter.Printf("✅ [Cleanup] Successfully removed temporary directories\n")
	}

	GinkgoWriter.Printf("✅ [Cleanup] Global test cleanup completed\n")
}
