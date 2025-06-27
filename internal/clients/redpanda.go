/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"

	"github.com/pcollinscode/provider-redpanda/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal redpanda credentials as JSON"
	keyClientId             = "client_id"
	keyClientSecret         = "client_secret"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}
		fmt.Fprintln(os.Stderr, "Beginning cred extract")
		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		fmt.Fprintln(os.Stderr, "Extracted creds")
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			fmt.Fprintln(os.Stderr, "We had an issue decoding the secret")
			fmt.Fprintf(os.Stderr, "%v", err)
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		fmt.Fprintf(os.Stderr, "Found data: %v", creds)

		// Set credentials in Terraform provider configuration.
		ps.Configuration = map[string]any{}
		if v, ok := creds[keyClientId]; ok {
			fmt.Fprintf(os.Stderr, "ClientId: %s", v)
			ps.Configuration[keyClientId] = v
		}
		fmt.Fprintln(os.Stderr, "Finished trying to get clientid")
		if v, ok := creds[keyClientSecret]; ok {
			ps.Configuration[keyClientSecret] = v
		}
		return ps, nil
	}
}
