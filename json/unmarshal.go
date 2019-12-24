package main

// findNested looks for a key named s in map m. If values in m map to other
// maps, findNested looks into them recursively. Returns true if found, and
// the value found.
func findNested(m map[string]interface{}, s string) (bool, interface{}) {
	// Try to find key s at this level
	for k, v := range m {
		if k == s {
			return true, v
		}
	}
	// Not found on this level, so try to find it nested
	for _, v := range m {
		nm, ok := v.(map[string]interface{})
		if ok {
			found, val := findNested(nm, s)
			if found {
				return found, val
			}
		}
	}
	// Not found recursively
	return false, nil
}