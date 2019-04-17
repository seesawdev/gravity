/*
Copyright 2018 Gravitational, Inc.

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
	"context"
	"path/filepath"
	"time"

	"github.com/gravitational/gravity/lib/constants"
	"github.com/gravitational/gravity/lib/storage"
	"github.com/gravitational/gravity/lib/storage/keyval"
	"github.com/gravitational/gravity/lib/utils"

	telecfg "github.com/gravitational/teleport/lib/config"
	"github.com/gravitational/teleport/lib/service"
	teleservices "github.com/gravitational/teleport/lib/services"
	teleutils "github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"
	"gopkg.in/check.v1"
)

type ProcessSuite struct{}

var _ = check.Suite(&ProcessSuite{})

func (s *ProcessSuite) TestAuthGatewayConfigReload(c *check.C) {
	// Initialize process with some default configuration.
	teleportConfig := telecfg.MakeSampleFileConfig()
	teleportConfig.DataDir = c.MkDir()
	teleportConfig.Proxy.CertFile = ""
	teleportConfig.Proxy.KeyFile = ""
	backend, err := keyval.NewBolt(keyval.BoltConfig{
		Path: filepath.Join(c.MkDir(), "test.db"),
	})
	c.Assert(err, check.IsNil)
	process := &Process{
		FieldLogger:       logrus.WithField(trace.Component, "process"),
		backend:           backend,
		tcfg:              *teleportConfig,
		authGatewayConfig: storage.DefaultAuthGateway(),
	}
	serviceConfig, err := process.buildTeleportConfig(process.authGatewayConfig)
	c.Assert(err, check.IsNil)
	process.TeleportProcess = &service.TeleportProcess{
		Supervisor: service.NewSupervisor("test"),
		Config:     serviceConfig,
	}

	// Update auth gateway setting that should trigger reload.
	process.reloadAuthGatewayConfig(storage.NewAuthGateway(
		storage.AuthGatewaySpecV1{
			ConnectionLimits: &storage.ConnectionLimits{
				MaxConnections: utils.Int64Ptr(50),
			},
		}))
	// Make sure reload event was broadcast.
	ch := make(chan service.Event)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	process.WaitForEvent(ctx, service.TeleportReloadEvent, ch)
	select {
	case <-ch:
	case <-ctx.Done():
		c.Fatal("didn't receive reload event")
	}

	// Now update principals.
	process.reloadAuthGatewayConfig(storage.NewAuthGateway(
		storage.AuthGatewaySpecV1{
			PublicAddr: &[]string{"example.com"},
		}))
	// Make sure process config is updated.
	config := process.TeleportProcess.Config
	comparePrincipals(c, config.Auth.PublicAddrs, []string{"example.com"})
	comparePrincipals(c, config.Proxy.SSHPublicAddrs, []string{"example.com"})
	comparePrincipals(c, config.Proxy.PublicAddrs, []string{"example.com"})
	comparePrincipals(c, config.Proxy.Kube.PublicAddrs, []string{"example.com"})
}

func comparePrincipals(c *check.C, addrs []teleutils.NetAddr, principals []string) {
	var hosts []string
	for _, addr := range addrs {
		hosts = append(hosts, addr.Host())
	}
	c.Assert(hosts, check.DeepEquals, principals)
}

func (s *ProcessSuite) TestClusterServices(c *check.C) {
	p := Process{
		context: context.TODO(),
	}

	// initially no services are running
	c.Assert(p.clusterServicesRunning(), check.Equals, false)

	service1Launched := make(chan bool)
	service1Done := make(chan bool)
	service1 := func(ctx context.Context) error {
		close(service1Launched)
		defer close(service1Done)
		select {
		case <-ctx.Done():
			return nil
		}
	}

	service2Launched := make(chan bool)
	service2Done := make(chan bool)
	service2 := func(ctx context.Context) error {
		close(service2Launched)
		defer close(service2Done)
		select {
		case <-ctx.Done():
			return nil
		}
	}

	// launch services
	err := p.startClusterServices([]clusterService{
		service1,
		service2,
	})
	c.Assert(err, check.IsNil)
	for i, ch := range []chan bool{service1Launched, service2Launched} {
		select {
		case <-ch:
		case <-time.After(time.Second):
			c.Fatalf("service%v wasn't launched", i+1)
		}
	}
	c.Assert(p.clusterServicesRunning(), check.Equals, true)

	// should not attempt to launch again
	err = p.startClusterServices([]clusterService{
		service1,
		service2,
	})
	c.Assert(err, check.NotNil)
	c.Assert(p.clusterServicesRunning(), check.Equals, true)

	// stop services
	err = p.stopClusterServices()
	c.Assert(err, check.IsNil)
	for i, ch := range []chan bool{service1Done, service2Done} {
		select {
		case <-ch:
		case <-time.After(time.Second):
			c.Fatalf("service%v wasn't stopped", i+1)
		}
	}
	c.Assert(p.clusterServicesRunning(), check.Equals, false)

	// should not attempt to stop again
	err = p.stopClusterServices()
	c.Assert(err, check.NotNil)
	c.Assert(p.clusterServicesRunning(), check.Equals, false)
}

func (s *ProcessSuite) TestReverseTunnelsFromTrustedClusters(c *check.C) {
	var testCases = []struct {
		clusters []teleservices.TrustedCluster
		tunnels  []telecfg.ReverseTunnel
		comment  string
	}{
		{
			clusters: nil,
			tunnels:  nil,
			comment:  "does nothing",
		},
		{
			clusters: []teleservices.TrustedCluster{
				storage.NewTrustedCluster("cluster1", storage.TrustedClusterSpecV2{
					Enabled:              false,
					ReverseTunnelAddress: "cluster1:3024",
					ProxyAddress:         "cluster1:443",
					Token:                "secret",
					Roles:                []string{constants.RoleAdmin},
				}),
			},
			tunnels: nil,
			comment: "ignores disabled clusters",
		},
		{
			clusters: []teleservices.TrustedCluster{
				storage.NewTrustedCluster("cluster1", storage.TrustedClusterSpecV2{
					Enabled:              true,
					ReverseTunnelAddress: "cluster1:3024",
					ProxyAddress:         "cluster1:443",
					Token:                "secret",
					Roles:                []string{constants.RoleAdmin},
				}),
				storage.NewTrustedCluster("cluster2", storage.TrustedClusterSpecV2{
					Enabled:              true,
					ReverseTunnelAddress: "cluster2:3024",
					ProxyAddress:         "cluster2:443",
					Token:                "secret",
					Roles:                []string{constants.RoleAdmin},
					Wizard:               true,
				}),
			},
			tunnels: []telecfg.ReverseTunnel{
				{
					DomainName: "cluster1",
					Addresses:  []string{"cluster1:3024"},
				},
				{
					DomainName: "cluster2",
					Addresses:  []string{"cluster2:3024"},
				},
			},
			comment: "considers all remote trusted clusters",
		},
	}
	for _, testCase := range testCases {
		backend, err := keyval.NewBolt(keyval.BoltConfig{
			Path: filepath.Join(c.MkDir(), "test.db"),
		})
		c.Assert(err, check.IsNil)
		for _, cluster := range testCase.clusters {
			_, err := backend.UpsertTrustedCluster(cluster)
			c.Assert(err, check.IsNil)
		}
		tunnels, err := reverseTunnelsFromTrustedClusters(backend)
		c.Assert(err, check.IsNil)
		c.Assert(tunnels, check.DeepEquals, testCase.tunnels, check.Commentf(testCase.comment))
	}
}
