// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	acl "github.com/pcollinscode/provider-redpanda/internal/controller/acl/acl"
	cluster "github.com/pcollinscode/provider-redpanda/internal/controller/cluster/cluster"
	providerconfig "github.com/pcollinscode/provider-redpanda/internal/controller/providerconfig"
	topic "github.com/pcollinscode/provider-redpanda/internal/controller/topic/topic"
	user "github.com/pcollinscode/provider-redpanda/internal/controller/user/user"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		acl.Setup,
		cluster.Setup,
		providerconfig.Setup,
		topic.Setup,
		user.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
