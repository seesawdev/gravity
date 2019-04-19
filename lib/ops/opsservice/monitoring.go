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

package opsservice

import (
	"context"

	"github.com/gravitational/gravity/lib/constants"
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/ops"
	"github.com/gravitational/gravity/lib/ops/events"
	"github.com/gravitational/gravity/lib/ops/monitoring"
	"github.com/gravitational/gravity/lib/storage"

	"github.com/gravitational/rigging"
	"github.com/gravitational/trace"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubelabels "k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// GetRetentionPolicies returns a list of retention policies for the site
func (o *Operator) GetRetentionPolicies(key ops.SiteKey) ([]monitoring.RetentionPolicy, error) {
	return o.cfg.Monitoring.GetRetentionPolicies()
}

// UpdateRetentionPolicy configures metrics retention policy
func (o *Operator) UpdateRetentionPolicy(req ops.UpdateRetentionPolicyRequest) error {
	err := req.Check()
	if err != nil {
		return trace.Wrap(err)
	}
	return o.cfg.Monitoring.UpdateRetentionPolicy(monitoring.RetentionPolicy{
		Name:     req.Name,
		Duration: req.Duration,
	})
}

// GetAlerts returns a list of configured monitoring alerts
func (o *Operator) GetAlerts(key ops.SiteKey) (alerts []storage.Alert, err error) {
	client, err := o.GetKubeClient()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	labels := kubelabels.Set{
		constants.MonitoringType: constants.MonitoringTypeAlert,
	}
	options := metav1.ListOptions{
		LabelSelector: labels.String(),
	}
	configmaps, err := client.Core().ConfigMaps(defaults.MonitoringNamespace).List(options)
	if err != nil {
		return nil, trace.Wrap(rigging.ConvertError(err))
	}

	var errors []error
	alerts = make([]storage.Alert, 0, len(configmaps.Items))
	for _, config := range configmaps.Items {
		data, ok := config.Data[constants.ResourceSpecKey]
		if !ok {
			continue
		}
		alert, err := storage.UnmarshalAlert([]byte(data))
		if err != nil {
			errors = append(errors, err)
			continue
		}
		alerts = append(alerts, alert)
	}

	if len(errors) != 0 {
		return nil, trace.NewAggregate(errors...)
	}

	return alerts, nil
}

// UpdateAlert updates the specified monitoring alert
func (o *Operator) UpdateAlert(ctx context.Context, key ops.SiteKey, alert storage.Alert) error {
	client, err := o.GetKubeClient()
	if err != nil {
		return trace.Wrap(err)
	}

	data, err := storage.MarshalAlert(alert)
	if err != nil {
		return trace.Wrap(err)
	}

	labels := map[string]string{
		constants.MonitoringType: constants.MonitoringTypeAlert,
	}
	err = updateConfigMap(client.Core().ConfigMaps(defaults.MonitoringNamespace),
		alert.GetName(), defaults.MonitoringNamespace, string(data), labels)
	if err != nil {
		return trace.Wrap(err)
	}

	events.Emit(ctx, o, events.ResourceAlertCreated, events.Fields{
		events.FieldName: alert.GetName(),
	})
	return nil
}

// DeleteAlert deletes the specified monitoring alert
func (o *Operator) DeleteAlert(ctx context.Context, key ops.SiteKey, name string) error {
	client, err := o.GetKubeClient()
	if err != nil {
		return trace.Wrap(err)
	}

	labels := kubelabels.Set{
		constants.MonitoringType: constants.MonitoringTypeAlert,
	}
	options := metav1.ListOptions{
		LabelSelector: labels.String(),
	}
	configmaps, err := client.Core().ConfigMaps(defaults.MonitoringNamespace).List(options)
	if err != nil {
		return trace.Wrap(rigging.ConvertError(err))
	}

	var alert *v1.ConfigMap
	for _, config := range configmaps.Items {
		if config.Name == name {
			alert = &config
			break
		}
	}
	if alert == nil {
		return trace.NotFound("alert %q not found", name)
	}

	err = client.Core().ConfigMaps(defaults.MonitoringNamespace).Delete(name, nil)
	if err != nil {
		return trace.Wrap(rigging.ConvertError(err))
	}

	events.Emit(ctx, o, events.ResourceAlertDeleted, events.Fields{
		events.FieldName: name,
	})
	return nil

}

// GetAlertTargets returns a list of configured monitoring alert targets
func (o *Operator) GetAlertTargets(key ops.SiteKey) (targets []storage.AlertTarget, err error) {
	client, err := o.GetKubeClient()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	data, err := getConfigMap(client.Core().ConfigMaps(defaults.MonitoringNamespace),
		constants.AlertTargetConfigMap)
	if err != nil {
		if trace.IsNotFound(err) {
			return nil, trace.NotFound("alert target not found")
		}
		return nil, trace.Wrap(err)
	}

	target, err := storage.UnmarshalAlertTarget([]byte(data))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return []storage.AlertTarget{target}, nil
}

// UpdateAlertTarget updates the cluster monitoring alert target
func (o *Operator) UpdateAlertTarget(ctx context.Context, key ops.SiteKey, target storage.AlertTarget) error {
	client, err := o.GetKubeClient()
	if err != nil {
		return trace.Wrap(err)
	}

	data, err := storage.MarshalAlertTarget(target)
	if err != nil {
		return trace.Wrap(err)
	}

	labels := map[string]string{
		constants.MonitoringType: constants.MonitoringTypeAlertTarget,
	}
	err = updateConfigMap(client.Core().ConfigMaps(defaults.MonitoringNamespace),
		constants.AlertTargetConfigMap, defaults.MonitoringNamespace, string(data), labels)
	if err != nil {
		return trace.Wrap(err)
	}

	events.Emit(ctx, o, events.ResourceAlertTargetCreated, events.Fields{})
	return nil

}

// DeleteAlertTarget deletes the cluster monitoring alert target
func (o *Operator) DeleteAlertTarget(ctx context.Context, key ops.SiteKey) error {
	client, err := o.GetKubeClient()
	if err != nil {
		return trace.Wrap(err)
	}

	err = rigging.ConvertError(client.Core().ConfigMaps(defaults.MonitoringNamespace).Delete(constants.AlertTargetConfigMap, nil))
	if err != nil {
		if trace.IsNotFound(err) {
			return trace.NotFound("no alert targets found")
		}
		return trace.Wrap(err)
	}

	events.Emit(ctx, o, events.ResourceAlertTargetDeleted, events.Fields{})
	return nil
}

func getConfigMap(client corev1.ConfigMapInterface, name string) (string, error) {
	config, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		return "", trace.Wrap(rigging.ConvertError(err))
	}

	data, ok := config.Data[constants.ResourceSpecKey]
	if !ok {
		return "", trace.NotFound("no resource found")
	}

	return data, nil
}

func updateConfigMap(client corev1.ConfigMapInterface, name, namespace, data string, labels map[string]string) error {
	config := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			constants.ResourceSpecKey: data,
		},
	}

	_, err := client.Create(config)
	err = rigging.ConvertError(err)
	if err == nil {
		return nil
	}

	if !trace.IsAlreadyExists(err) {
		return trace.Wrap(err)
	}

	_, err = client.Update(config)
	return trace.Wrap(rigging.ConvertError(err))
}
