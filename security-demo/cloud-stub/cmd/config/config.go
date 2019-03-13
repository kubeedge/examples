package config

import (
	"io/ioutil"
	"fmt"
	"github.com/hashicorp/hcl"
)

// CloudStubConfig is HCL config data
type CloudStubConfig struct {
	PlacementURL  string `hcl:"placementURL"`
}

// ParseConfig parses the given HCL file into a SidecarConfig struct
func ParseConfig(file string) (stubConfig *CloudStubConfig, err error) {
	stubConfig = &CloudStubConfig{}

	// Read HCL file
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	hclText := string(dat)

	// Parse HCL
	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&stubConfig, hclParseTree); err != nil {
		return nil, err
	}

	fmt.Println("config file : %s, config : %s", file, stubConfig.PlacementURL)
	return stubConfig, nil
}

