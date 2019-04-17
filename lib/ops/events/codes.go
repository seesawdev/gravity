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

package events

import "github.com/gravitational/teleport/lib/events"

var (
	// OperationInstallStart is emitted when a cluster installation starts.
	OperationInstallStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationInstallStartCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} is being installed",
	}
	// OperationInstallComplete is emitted when a cluster installation successfully completes.
	OperationInstallComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationInstallCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} has been installed",
	}
	// OperationInstallFailure is emitted when a cluster installation fails.
	OperationInstallFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationInstallFailureCode,
		Severity: events.SeverityError,
		Message:  "Cluster {{.cluster}} install has failed",
	}
	// OperationExpandStart is emitted when a new node starts joining the cluster.
	OperationExpandStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationExpandStartCode,
		Severity: events.SeverityInfo,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} is joining the cluster {{.cluster}}",
	}
	// OperationExpandComplete is emitted when a node has successfully joined the cluster.
	OperationExpandComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationExpandCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} has joined the cluster {{.cluster}}",
	}
	// OperationExpandFailure is emitted when a node fails to join the cluster.
	OperationExpandFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationExpandFailureCode,
		Severity: events.SeverityError,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} has failed to join the cluster {{.cluster}}",
	}
	// OperationShrinkStart is emitted when a node is leaving the cluster.
	OperationShrinkStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationShrinkStartCode,
		Severity: events.SeverityInfo,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} is leaving the cluster {{.cluster}}",
	}
	// OperationShrinkComplete is emitted when a node has left the cluster.
	OperationShrinkComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationShrinkCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} has left the cluster {{.cluster}}",
	}
	// OperationShrinkFailure is emitted when a node fails to leave the cluster.
	OperationShrinkFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationShrinkFailureCode,
		Severity: events.SeverityError,
		Message:  "Node {{.hostname}}/{{.ip}}/{{.role}} has failed to leave the cluster {{.cluster}}",
	}
	// OperationUpdateStart is emitted when cluster upgrade is started.
	OperationUpdateStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationUpdateStartCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} is upgrading to version {{.version}}",
	}
	// OperationUpdateCompete is emitted when cluster upgrade successfully finishes.
	OperationUpdateComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationUpdateCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} has been upgraded to version {{.version}}",
	}
	// OperationUpdateFailure is emitted when cluster upgrade fails.
	OperationUpdateFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationUpdateFailureCode,
		Severity: events.SeverityError,
		Message:  "Cluster {{.cluster}} upgrade to version {{.version}} has failed",
	}
	// OperationUninstallStart is emitted when cluster uninstall is launched.
	OperationUninstallStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationUninstallStartCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} is being uninstalled",
	}
	// OperationUninstallComplete is emitted when cluster has been uninstalled.
	OperationUninstallComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationUninstallCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} has been uninstalled",
	}
	// OperationUninstallFailure is emitted when cluster uninstall fails.
	OperationUninstallFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationUninstallFailureCode,
		Severity: events.SeverityError,
		Message:  "Cluster {{.cluster}} uninstall has failed",
	}
	// OperationGCStart is emitted when garbage collection is started on a cluster.
	OperationGCStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationGCStartCode,
		Severity: events.SeverityInfo,
		Message:  "Running garbage collection on cluster {{.cluster}}",
	}
	// OperationGCComplete is emitted when cluster garbage collection successfully completes.
	OperationGCComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationGCCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Garbage collection on cluster {{.cluster}} has finished",
	}
	// OperationGCFailure is emitted when cluster garbage collection fails.
	OperationGCFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationGCFailureCode,
		Severity: events.SeverityError,
		Message:  "Garbage collection on cluster {{.cluster}} has failed",
	}
	// OperationEnvStart is emitted when cluster runtime environment update is launched.
	OperationEnvStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationEnvStartCode,
		Severity: events.SeverityInfo,
		Message:  "Updating runtime environment on cluster {{.cluster}}",
	}
	// OperationEnvComplete is emitted when cluster runtime environment update successfully completes.
	OperationEnvComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationEnvCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Runtime environment on cluster {{.cluster}} has been updated",
	}
	// OperationEnvFailure is emitted when cluster runtime environment update fails.
	OperationEnvFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationEnvFailureCode,
		Severity: events.SeverityError,
		Message:  "Failed to update runtime environment on cluster {{.cluster}}",
	}
	// OperationConfigStart is emitted when cluster configuration update launches.
	OperationConfigStart = events.Event{
		Name:     OperationStarted,
		Code:     OperationConfigStartCode,
		Severity: events.SeverityInfo,
		Message:  "Updating cluster {{.cluster}} configuration",
	}
	// OperationConfigComplete is emitted when cluster configuration update successfully completes.
	OperationConfigComplete = events.Event{
		Name:     OperationCompleted,
		Code:     OperationConfigCompleteCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} configuration has been updated",
	}
	// OperationConfigFailure is emitted when cluster configuration update fails.
	OperationConfigFailure = events.Event{
		Name:     OperationFailed,
		Code:     OperationConfigFailureCode,
		Severity: events.SeverityError,
		Message:  "Failed to update cluster {{.cluster}} configuration",
	}
	// ResourceUserCreated is emitted when a user is created/updated.
	ResourceUserCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceUserCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created user {{.name}}",
	}
	// ResourceUserDeleted is emitted when a user is deleted.
	ResourceUserDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceUserDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted user {{.name}}",
	}
	// ResourceTokenCreated is emitted when a token is created/updated.
	ResourceTokenCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceTokenCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created token for user {{.owner}}",
	}
	// ResourceTokenDeleted is emitted when a token is deleted.
	ResourceTokenDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceTokenDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted token for user {{.owner}}",
	}
	// ResourceGithubConnectorCreated is emitted when a Github connector is created/updated.
	ResourceGithubConnectorCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceGithubConnectorCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created Github connector {{.name}}",
	}
	// ResourceGithubConnectorDeleted is emitted when a Github connector is deleted.
	ResourceGithubConnectorDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceGithubConnectorDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted Github connector {{.name}}",
	}
	// ResourceLogForwarderCreated is emitted when a log forwarder is created/updated.
	ResourceLogForwarderCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceLogForwarderCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created log forwarder {{.name}}",
	}
	// ResourceLogForwarderDeleted is emitted when a log forwarder is deleted.
	ResourceLogForwarderDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceLogForwarderDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted log forwarder {{.name}}",
	}
	// ResourceTLSKeyPairCreated is emitted when cluster web certificate is updated.
	ResourceTLSKeyPairCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceTLSKeyPairCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} installed cluster web certificate",
	}
	// ResourceTLSKeyPairDeleted is emitted when cluster web certificate is deleted.
	ResourceTLSKeyPairDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceLogForwarderDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted cluster web certificate",
	}
	// ResourceAuthPreferenceCreated is emitted when cluster auth preference is updated.
	ResourceAuthPreferenceCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceAuthPreferenceCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} updated cluster authentication preference",
	}
	// ResourceSMTPConfigCreated is emitted when SMTP configuration is created/updated.
	ResourceSMTPConfigCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceSMTPConfigCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} updated cluster SMTP configuration",
	}
	// ResourceSMTPConfigDeleted is emitted when SMTP configuration is deleted.
	ResourceSMTPConfigDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceSMTPConfigDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted cluster SMTP configuration",
	}
	// ResourceAlertCreated is emitted when monitoring alert is created/updated.
	ResourceAlertCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceAlertCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created monitoring alert {{.name}}",
	}
	// ResourceAlertDeleted is emitted when monitoring alert is deleted.
	ResourceAlertDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceAlertDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted monitoring alert {{.name}}",
	}
	// ResourceAlertTargetCreated is emitted when monitoring alert target is created/updated.
	ResourceAlertTargetCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceAlertTargetCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} created monitoring alert target {{.name}}",
	}
	// ResourceAlertTargetDeleted is emitted when monitoring alert target is deleted.
	ResourceAlertTargetDeleted = events.Event{
		Name:     ResourceDeleted,
		Code:     ResourceAlertTargetDeletedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} deleted monitoring alert target {{.name}}",
	}
	// ResourceAuthGatewayCreated is emitted when cluster auth gateway settings are updated.
	ResourceAuthGatewayCreated = events.Event{
		Name:     ResourceCreated,
		Code:     ResourceAuthGatewayCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} updated cluster authentication gateway settings",
	}
	// UserInviteCreated is emitted when a user invite is created.
	UserInviteCreated = events.Event{
		Name:     InviteCreated,
		Code:     UserInviteCreatedCode,
		Severity: events.SeverityInfo,
		Message:  "User {{.user}} invited user {{.name}} with roles {{.roles}}",
	}
	// ClusterUnhealthy is emitted when cluster becomes unhealthy.
	ClusterUnhealthy = events.Event{
		Name:     ClusterDegraded,
		Code:     ClusterUnhealthyCode,
		Severity: events.SeverityWarning,
		Message:  "Cluster {{.cluster}} is degraded: {{.reason}}",
	}
	// ClusterHealthy is emitted when cluster becomes healthy.
	ClusterHealthy = events.Event{
		Name:     ClusterActivated,
		Code:     ClusterHealthyCode,
		Severity: events.SeverityInfo,
		Message:  "Cluster {{.cluster}} has become healthy",
	}
)

var (
	OperationInstallStartCode          = "G0001I"
	OperationInstallCompleteCode       = "G0002I"
	OperationInstallFailureCode        = "G0002E"
	OperationExpandStartCode           = "G0003I"
	OperationExpandCompleteCode        = "G0004I"
	OperationExpandFailureCode         = "G0004E"
	OperationShrinkStartCode           = "G0005I"
	OperationShrinkCompleteCode        = "G0006I"
	OperationShrinkFailureCode         = "G0006E"
	OperationUpdateStartCode           = "G0007I"
	OperationUpdateCompleteCode        = "G0008I"
	OperationUpdateFailureCode         = "G0008E"
	OperationUninstallStartCode        = "G0009I"
	OperationUninstallCompleteCode     = "G0010I"
	OperationUninstallFailureCode      = "G0010E"
	OperationGCStartCode               = "G0011I"
	OperationGCCompleteCode            = "G0012I"
	OperationGCFailureCode             = "G0012E"
	OperationEnvStartCode              = "G0013I"
	OperationEnvCompleteCode           = "G0014I"
	OperationEnvFailureCode            = "G0014E"
	OperationConfigStartCode           = "G0015I"
	OperationConfigCompleteCode        = "G0016I"
	OperationConfigFailureCode         = "G0016E"
	ResourceUserCreatedCode            = "G1000I"
	ResourceUserDeletedCode            = "G2000I"
	ResourceTokenCreatedCode           = "G1001I"
	ResourceTokenDeletedCode           = "G2001I"
	ResourceGithubConnectorCreatedCode = "G1002I"
	ResourceGithubConnectorDeletedCode = "G2002I"
	ResourceLogForwarderCreatedCode    = "G1003I"
	ResourceLogForwarderDeletedCode    = "G2003I"
	ResourceTLSKeyPairCreatedCode      = "G1004I"
	ResourceTLSKeyPairDeletedCode      = "G2004I"
	ResourceAuthPreferenceCreatedCode  = "G1005I"
	ResourceSMTPConfigCreatedCode      = "G1006I"
	ResourceSMTPConfigDeletedCode      = "G2006I"
	ResourceAlertCreatedCode           = "G1007I"
	ResourceAlertDeletedCode           = "G2007I"
	ResourceAlertTargetCreatedCode     = "G1008I"
	ResourceAlertTargetDeletedCode     = "G2008I"
	ResourceAuthGatewayCreatedCode     = "G1009I"
	UserInviteCreatedCode              = "G1010I"
	ClusterUnhealthyCode               = "G3000W"
	ClusterHealthyCode                 = "G3001I"
)
