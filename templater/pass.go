// This file contains the implementation of a template function for retrieving variables from 'pass', the standard UNIX
// password manager.
package templater

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/polydawn/meep"
)

type PassError struct {
	meep.TraitAutodescribing
	meep.TraitCausable
	Output string
}

func GetFromPass(key string) (string, error) {
	fmt.Fprintf(os.Stderr, "Attempting to look up %s in pass\n", key)
	pass := exec.Command("pass", "show", key)

	output, err := pass.CombinedOutput()
	if err != nil {
		return "", meep.New(
			&PassError{Output: string(output)},
			meep.Cause(err),
		)
	}

	return string(output), nil
}
