package templater

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/polydawn/meep"
	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/util"
)

func GetResource(c *context.Context, kind string, name string, namespace string) (string, error) {
	if namespace == "" {
		namespace = "default"
	}

	fmt.Fprintf(os.Stderr, "Attempting to retrieve last applied configuration for %s %s in namespace %s",
		kind, name, namespace)

	return getLastAppliedConfiguration(c, kind, name, namespace)
}

// Retrieves the last applied configuration of a resource. Kontemplate always uses the correct flags to store the
// necessary annotation, but this may not work with resources created using 'kubectl create' by other tools.
// The call requires at least kubectl v1.6.0.
func getLastAppliedConfiguration(c *context.Context, kind string, name string, namespace string) (string, error) {
	args := []string{
		fmt.Sprintf("--context=%s", c.Name),
		fmt.Sprintf("--namespace=%s", namespace),
		"apply", "view-last-applied", kind, name,
	}

	kubectl := exec.Command("kubectl", args...)

	w := bytes.NewBuffer([]byte{})
	kubectl.Stdout = w

	e := bytes.NewBuffer([]byte{})
	kubectl.Stderr = e

	if err := kubectl.Start(); err != nil {
		return "", meep.New(&util.KubeCtlError{
			Args:   args,
			Stderr: e.String(),
		}, meep.Cause(err))
	}

	if err := kubectl.Wait(); err != nil {
		return "", meep.New(&util.KubeCtlError{
			Args:   args,
			Stderr: e.String(),
		}, meep.Cause(err))
	}

	return w.String(), nilsu
}
