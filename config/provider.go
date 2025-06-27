/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/pkg/config"

	acl "github.com/pcollinscode/provider-redpanda/config/acl"
	cluster "github.com/pcollinscode/provider-redpanda/config/cluster"
	topic "github.com/pcollinscode/provider-redpanda/config/topic"
	user "github.com/pcollinscode/provider-redpanda/config/user"
)

const (
	resourcePrefix = "redpanda"
	modulePath     = "github.com/pcollinscode/provider-redpanda"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("pcollinscode.com"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		acl.Configure,
		cluster.Configure,
		topic.Configure,
		user.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
