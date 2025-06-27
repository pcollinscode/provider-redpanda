package repository

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("redpanda_cluster", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "github"
		r.ShortGroup = "cluster"

		delete(r.TerraformResource.Schema, "aws_private_link")
		delete(r.TerraformResource.Schema, "azure_private_link")
		delete(r.TerraformResource.Schema, "customer_managed_resources")
		delete(r.TerraformResource.Schema, "gcp_private_service_connect")
		delete(r.TerraformResource.Schema, "http_proxy")
		delete(r.TerraformResource.Schema, "kafka_api")
		delete(r.TerraformResource.Schema, "kafka_connect")
		delete(r.TerraformResource.Schema, "maintenance_window_config")
		delete(r.TerraformResource.Schema, "schema_registry")
		delete(r.TerraformResource.Schema, "prometheus")
		delete(r.TerraformResource.Schema, "redpanda_console")
		delete(r.TerraformResource.Schema, "state_description")
	})
}
