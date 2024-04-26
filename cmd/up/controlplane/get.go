// Copyright 2022 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	spacesv1beta1 "github.com/upbound/up-sdk-go/apis/spaces/v1beta1"
	cp "github.com/upbound/up-sdk-go/service/controlplanes"
	"github.com/upbound/up/internal/upbound"
	"github.com/upbound/up/internal/upterm"
)

// AfterApply sets default values in command after assignment and validation.
func (c *getCmd) AfterApply(kongCtx *kong.Context, upCtx *upbound.Context) error {
	kongCtx.Bind(pterm.DefaultTable.WithWriter(kongCtx.Stdout).WithSeparator("   "))
	return nil
}

// getCmd gets a single control plane in an account on Upbound.
type getCmd struct {
	Name string `arg:"" required:"" help:"Name of control plane." predictor:"ctps"`
}

// Run executes the get command.
func (c *getCmd) Run(ctx context.Context, printer upterm.ObjectPrinter, p pterm.TextPrinter, upCtx *upbound.Context, client client.Client) error {
	var ctp spacesv1beta1.ControlPlane
	ns, _, err := upCtx.Kubecfg.Namespace()
	if err != nil {
		return errors.Wrap(err, "error getting namespace")
	}

	if err := client.Get(ctx, types.NamespacedName{Name: c.Name, Namespace: ns}, &ctp); err != nil {
		if kerrors.IsNotFound(err) {
			p.Printfln("Control plane %s not found", c.Name)
			return nil
		}

		return errors.Wrap(err, "error getting control plane")
	}

	return tabularPrint(ctp, printer, upCtx)
}

// EmptyControlPlaneConfiguration returns an empty ControlPlaneConfiguration with default values.
func EmptyControlPlaneConfiguration() cp.ControlPlaneConfiguration {
	configuration := cp.ControlPlaneConfiguration{}
	configuration.Status = cp.ConfigurationInstallationQueued
	return configuration
}
