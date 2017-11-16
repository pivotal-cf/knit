package main_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	patcher          string
	cfReleaseRepo    string
	cfPatchesDir     string
	diegoReleaseRepo string
	diegoPatchesDir  string
)

var _ = BeforeSuite(func() {
	var err error
	patcher, err = gexec.Build("github.com/pivotal-cf/knit")
	Expect(err).NotTo(HaveOccurred())

	cfPatchesDir = os.Getenv("CF_PATCHES_DIR")

	diegoPatchesDir = os.Getenv("DIEGO_PATCHES_DIR")

	if os.Getenv("CF_RELEASE_DIR") == "" {
		Fail("CF_RELEASE_DIR is a required env var")
	}

	if cfPatchesDir == "" {
		Fail("CF_PATCHES_DIR is a required env var")
	}

	if os.Getenv("DIEGO_RELEASE_DIR") == "" {
		Fail("DIEGO_RELEASE_DIR is a required env var")
	}

	if diegoPatchesDir == "" {
		Fail("DIEGO_PATCHES_DIR is a required env var")
	}
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestApplyPatches(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApplyPatches Suite")
}
