package testutils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

// DetectLeakyGoroutines is a helper function that detects and prevents leaked goroutines during testing.
// It uses the gleak library to track goroutine leaks and asserts that no goroutines are leaked.
// This function should be called at the beginning of each test case that involves goroutines.
// It uses ginkgo's DeferCleanup function to ensure that the cleanup code is executed after each test case.
func DetectLeakyGoroutines() {
	// Capture the initial set of goroutines
	nonLeakyGoroutines := Goroutines()

	// Register a cleanup function to assert that no goroutines are leaked
	DeferCleanup(func() {
		// Assert that no goroutines are leaked
		Eventually(Goroutines).ShouldNot(HaveLeaked(nonLeakyGoroutines))
	})
}
