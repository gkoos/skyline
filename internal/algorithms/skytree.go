package algorithms

import "github.com/gkoos/skyline/types"

var defaultSkyTreeConfig = types.SkyTreeConfig{}

// SkyTree is a dummy implementation of SkyTree skyline algorithm
func SkyTree(data types.Dataset, prefs types.Preference, cfg *types.SkyTreeConfig) types.Dataset {
	if cfg == nil {
		cfg = &defaultSkyTreeConfig
	}
	
	// TODO: implement real algorithm
	return data[:1] // dummy: return first point only
}
