/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/ops/opsservice"
	"github.com/gravitational/gravity/lib/processconfig"
	"github.com/gravitational/gravity/lib/storage"

	"github.com/gravitational/teleport/lib/config"
	"github.com/gravitational/teleport/lib/service"

	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"
)

// buildTeleportConfig builds configuration object that will be used to
// start embedded Teleport services.
func (p *Process) buildTeleportConfig(authGatewayConfig storage.AuthGateway) (*service.Config, error) {
	configFromImport, err := p.getTeleportConfigFromImportState()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	// If we're running inside Kubernetes, Teleport configuration is stored
	// in a package that's available in what we call "import state".
	fileConfig := p.tcfg
	if configFromImport != nil {
		fileConfig = *configFromImport
	}
	err = processconfig.MergeTeleConfigFromEnv(&fileConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	// Apply user-defined configuration on top of the file config. This is
	// what users configure via AuthGateway resource.
	if authGatewayConfig != nil {
		authGatewayConfig.ApplyToTeleportConfig(&fileConfig)
	}
	fileConfig.Auth.ReverseTunnels, err = reverseTunnelsFromTrustedClusters(p.backend)
	if err != nil && !trace.IsNotFound(err) {
		return nil, trace.Wrap(err)
	}
	serviceConfig := service.MakeDefaultConfig()
	err = config.ApplyFileConfig(&fileConfig, serviceConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	serviceConfig.Auth.StorageConfig.Params["path"] = fileConfig.DataDir
	if len(serviceConfig.AuthServers) == 0 && serviceConfig.Auth.Enabled {
		serviceConfig.AuthServers = append(serviceConfig.AuthServers, serviceConfig.Auth.SSHAddr)
	}
	// Teleport will be using Gravity backend implementation.
	serviceConfig.Identity = p.identity
	serviceConfig.Trust = p.identity
	serviceConfig.Presence = p.backend
	serviceConfig.Provisioner = p.identity
	serviceConfig.Proxy.DisableWebInterface = true
	serviceConfig.Proxy.DisableWebService = true
	serviceConfig.Access = p.identity
	serviceConfig.Console = logrus.StandardLogger().Writer()
	serviceConfig.ClusterConfiguration = p.identity
	return serviceConfig, nil
}

// getOrInitAuthGatewayConfig returns auth gateway configuration.
//
// If it's not found, it's first initialized with default values.
func (p *Process) getOrInitAuthGatewayConfig() (storage.AuthGateway, error) {
	if !p.inKubernetes() {
		// We're not running inside Kubernetes, so this is likely an installer
		// process which doesn't support auth gateway reconfiguration.
		return nil, nil
	}
	cluster, err := p.backend.GetLocalSite(defaults.SystemAccountID)
	if err != nil {
		if trace.IsNotFound(err) {
			// There's no local cluster which likely means that process is
			// started in "init" mode (which imports cluster packages and
			// backend data from the "import" state) so auth gateway will
			// be initialized during the actual process startup.
			return nil, nil
		}
		return nil, trace.Wrap(err)
	}

	client, err := tryGetPrivilegedKubeClient()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	authGateway, err := opsservice.GetAuthGateway(client, p.identity)
	if err != nil && !trace.IsNotFound(err) {
		return nil, trace.Wrap(err)
	}
	if authGateway != nil {
		p.Debug("Auth gateway config resource is already initialized.")
		return authGateway, nil
	}
	// Auth gateway resource does not exist, initialize it with default
	// values taken from Teleport config.
	p.Info("Initializing auth gateway config resource.")
	authPreference, err := p.getAuthPreference()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	authGateway = storage.DefaultAuthGateway()
	err = authGateway.SetAuthPreference(authPreference)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	// Initially the local cluster name is set as a principal.
	authGateway.SetSSHPublicAddrs([]string{cluster.Domain})
	authGateway.SetKubernetesPublicAddrs([]string{cluster.Domain})
	err = opsservice.UpsertAuthGateway(client, p.identity, authGateway)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return authGateway, nil
}

func (p *Process) getAuthGatewayConfig() (storage.AuthGateway, error) {
	client, err := tryGetPrivilegedKubeClient()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return opsservice.GetAuthGateway(client, p.identity)
}
