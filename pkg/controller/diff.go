package controller

import (
	"encoding/json"
	"fmt"
)

type Drift struct {
	// Indicate if there are any changes in live and desired state
	HasChanges bool

	// count of missing resources which are present in desired but not in live
	Missing []ResourceState

	// count of extra resources which are present in live but not in desired
	Extra []ResourceState

	// Resources present in both but with different specs
	Modified []ResourceState

	// summary of the drift
	Summary string
}

// dirft detection
func Compare(liveState []ResourceState, desiredState []ResourceState) (Drift, error) {

	livemap := make(map[string]ResourceState)
	desiredmap := make(map[string]ResourceState)

	for _, resource := range liveState {
		key := fmt.Sprintf("%s/%s/%s", resource.Kind, resource.Name, resource.Namespace)
		livemap[key] = resource
	}

	for _, resource := range desiredState {
		key := fmt.Sprintf("%s/%s/%s", resource.Kind, resource.Name, resource.Namespace)
		desiredmap[key] = resource
	}

	var missing []ResourceState
	for _, resource := range desiredState {
		key := fmt.Sprintf("%s/%s/%s", resource.Kind, resource.Name, resource.Namespace)
		if _, ok := livemap[key]; !ok {
			missing = append(missing, resource)
		}
	}

	var extra []ResourceState
	var modified []ResourceState

	for _, resource := range liveState {
		key := fmt.Sprintf("%s/%s/%s", resource.Kind, resource.Name, resource.Namespace)

		// TODO: Refactor to be more flexible for all system resources
		if resource.Name == "kubernetes" && resource.Kind == "Service" {
			continue
		}

		if desired, ok := desiredmap[key]; !ok {
			extra = append(extra, resource)
		} else {
			// compare specs(same resources in both)
			if !specsEqual(resource.Spec, desired.Spec) {
				modified = append(modified, desired)
				
				// Debug: Log the difference (can remove later)
				liveJSON, _ := json.Marshal(resource.Spec)
				desiredJSON, _ := json.Marshal(desired.Spec)
				fmt.Printf("DEBUG: Spec mismatch for %s/%s/%s\n", resource.Kind, resource.Name, resource.Namespace)
				fmt.Printf("  Live: %s\n", string(liveJSON))
				fmt.Printf("  Desired: %s\n", string(desiredJSON))
			}

			
		}
	}

	summary := fmt.Sprintf("Missing: %d, Extra: %d, Modified: %d", len(missing), len(extra), len(modified))
	if len(missing) == 0 && len(extra) == 0 && len(modified) == 0 {
		summary = "Synced - no changes detected"
	}

	return Drift{
		HasChanges: len(missing) > 0 || len(extra) > 0 || len(modified) > 0,
		Missing:    missing,
		Extra:      extra,
		Modified:   modified,
		Summary:    summary,
	}, nil
}

// compares two specs
func specsEqual(spec1, spec2 any) bool {
	if spec1 == nil && spec2 == nil {
		return true
	}
	if spec1 == nil || spec2 == nil {
		return false
	}

	json1, err1 := json.Marshal(spec1)
	json2, err2 := json.Marshal(spec2)

	if err1 != nil || err2 != nil {
		return false
	}

	return string(json1) == string(json2)
}
