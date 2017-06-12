package util

import "github.com/polydawn/meep"

// Common error type for kubectl calls
type KubeCtlError struct {
	meep.TraitAutodescribing
	meep.TraitCausable

	// The argument list that kubectl was called with
	Args []string

	// Kubectl stderr (if present)
	Stderr string
}

// Merges two maps together. Values from the second map override values in the first map.
// The returned map is new if anything was changed.
func Merge(in1 *map[string]interface{}, in2 *map[string]interface{}) *map[string]interface{} {
	if in1 == nil || len(*in1) == 0 {
		return in2
	}

	if in2 == nil || len(*in2) == 0 {
		return in1
	}

	new := make(map[string]interface{})
	for k, v := range *in1 {
		new[k] = v
	}

	for k, v := range *in2 {
		new[k] = v
	}

	return &new
}
