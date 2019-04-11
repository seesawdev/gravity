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

package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gravitational/gravity/lib/constants"
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/fsm"
	libfsm "github.com/gravitational/gravity/lib/fsm"
	"github.com/gravitational/gravity/lib/localenv"
	"github.com/gravitational/gravity/lib/ops"
	"github.com/gravitational/gravity/lib/storage"
	"github.com/gravitational/gravity/lib/update"

	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"
)

func newUpdater(ctx context.Context, localEnv, updateEnv *localenv.LocalEnvironment, init updateInitializer) (*update.Updater, error) {
	teleportClient, err := localEnv.TeleportClient(constants.Localhost)
	if err != nil {
		return nil, trace.Wrap(err, "failed to create a teleport client")
	}
	proxy, err := teleportClient.ConnectToProxy(ctx)
	if err != nil {
		return nil, trace.Wrap(err, "failed to connect to teleport proxy")
	}
	clusterEnv, err := localEnv.NewClusterEnvironment()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	if clusterEnv.Client == nil {
		return nil, trace.BadParameter("this operation can only be executed on one of the master nodes")
	}
	operator := clusterEnv.Operator
	cluster, err := operator.GetLocalSite()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = init.validatePreconditions(localEnv, operator, *cluster)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	key, err := init.newOperation(operator, *cluster)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	logger := logrus.WithField("operation", key)
	defer func() {
		r := recover()
		panicked := r != nil
		if err != nil || panicked {
			logger.WithError(err).Warn("Operation failed.")
			var msg string
			if err != nil {
				msg = err.Error()
			}
			if errReset := ops.FailOperationAndResetCluster(*key, operator, msg); errReset != nil {
				logrus.WithFields(logrus.Fields{
					logrus.ErrorKey: errReset,
					"operation":     key,
				}).Warn("Failed to mark operation as failed.")
			}
		}
		if r != nil {
			panic(r)
		}
	}()
	operation, err := operator.GetSiteOperation(*key)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	// Create the operation plan so it can be replicated on remote nodes
	plan, err := init.newOperationPlan(ctx, operator, *cluster, *operation, localEnv, updateEnv, clusterEnv)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = update.SyncOperationPlan(clusterEnv.Backend, updateEnv.Backend, *plan,
		(storage.SiteOperation)(*operation))
	if err != nil {
		return nil, trace.Wrap(err)
	}
	req := init.updateDeployRequest(deployAgentsRequest{
		clusterState: cluster.ClusterState,
		clusterName:  cluster.Domain,
		clusterEnv:   clusterEnv,
		proxy:        proxy,
		nodeParams:   constants.RPCAgentSyncPlanFunction,
	})
	deployCtx, cancel := context.WithTimeout(ctx, defaults.AgentDeployTimeout)
	defer cancel()
	logger.WithField("request", req).Debug("Deploying agents on nodes.")
	localEnv.Println("Deploying agents on nodes")
	creds, err := deployAgents(deployCtx, req)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	if bool(localEnv.Silent) {
		// FIXME: keep the legacy behavior of reporting the operation ID in quiet mode.
		// This is still used by robotest to fetch the operation ID
		fmt.Println(key.OperationID)
	}
	runner := libfsm.NewAgentRunner(creds)
	updater, err := init.newUpdater(ctx, operator, *operation, localEnv, updateEnv, clusterEnv, runner)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return updater, nil
}

type updateInitializer interface {
	validatePreconditions(localEnv *localenv.LocalEnvironment, operator ops.Operator, cluster ops.Site) error
	newOperation(ops.Operator, ops.Site) (*ops.SiteOperationKey, error)
	newOperationPlan(ctx context.Context,
		operator ops.Operator,
		cluster ops.Site,
		operation ops.SiteOperation,
		localEnv, updateEnv *localenv.LocalEnvironment,
		clusterEnv *localenv.ClusterEnvironment) (*storage.OperationPlan, error)
	newUpdater(ctx context.Context,
		operator ops.Operator,
		operation ops.SiteOperation,
		localEnv, updateEnv *localenv.LocalEnvironment,
		clusterEnv *localenv.ClusterEnvironment,
		runner fsm.AgentRepository,
	) (*update.Updater, error)
	updateDeployRequest(deployAgentsRequest) deployAgentsRequest
}

type updater interface {
	io.Closer
	Run(ctx context.Context) error
	RunPhase(ctx context.Context, phase string, phaseTimeout time.Duration, force bool) error
	RollbackPhase(ctx context.Context, phase string, phaseTimeout time.Duration, force bool) error
	Complete(error) error
}
