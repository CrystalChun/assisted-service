/*
Copyright 2020.

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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var agentclusterinstalllog = logf.Log.WithName("agentclusterinstall-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *AgentClusterInstall) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-hiveextension-openshift-io-v1beta1-agentclusterinstall,mutating=true,failurePolicy=fail,sideEffects=None,groups=hiveextension.openshift.io,resources=agentclusterinstalls,verbs=create;update,versions=v1beta1,name=magentclusterinstall.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &AgentClusterInstall{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *AgentClusterInstall) Default() {
	agentclusterinstalllog.Info("default", "name", r.Name)

}

//+kubebuilder:webhook:path=/validate-hiveextension-openshift-io-v1beta1-agentclusterinstall,mutating=false,failurePolicy=fail,sideEffects=None,groups=hiveextension.openshift.io,resources=agentclusterinstalls,verbs=create;update,versions=v1beta1,name=vagentclusterinstall.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AgentClusterInstall{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateCreate() (admission.Warnings, error) {
	agentclusterinstalllog.Info("validate create", "name", r.Name)

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	agentclusterinstalllog.Info("validate update", "name", r.Name)

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateDelete() (admission.Warnings, error) {
	agentclusterinstalllog.Info("validate delete", "name", r.Name)

	return nil, nil
}
