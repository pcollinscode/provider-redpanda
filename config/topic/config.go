package repository

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("redpanda_topic", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "github"
		r.ShortGroup = "topic"

		r.References["cluster"] = config.Reference{
			Type: "github.com/pcollinscode/provider-redpanda/apis/repository/v1alpha1.Cluster",
		}
	})
}
