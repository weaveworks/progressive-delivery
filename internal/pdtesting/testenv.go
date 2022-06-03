package pdtesting

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/weaveworks/weave-gitops/pkg/testutils"
)

func getRepoRoot() string {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(cmdOut))
}

func CreateTestEnv() (*testutils.K8sTestEnv, error) {
	envTestPath := fmt.Sprintf("%s/tools/bin/envtest", getRepoRoot())
	os.Setenv("KUBEBUILDER_ASSETS", envTestPath)

	var err error
	k8sEnv, err := testutils.StartK8sTestEnvironment([]string{
		"../../tools/testcrds",
	})

	if err != nil {
		return nil, err
	}

	return k8sEnv, nil
}
