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
	"fmt"
	"strings"

	"github.com/openshift/assisted-service/models"

	"github.com/go-openapi/swag"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	mutableFields = []string{"ClusterMetadata", "IgnitionEndpoint"}
	// log is for logging in this package.
	agentclusterinstalllog = logf.Log.WithName("agentclusterinstall-resource")
)

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

	//TODO When is usermanagednetworking required??

	// if UserNetworkManagement is not set by the user apply the defaults:
	// true for SNO and false for multi-node
	if !installAlreadyStarted(r.Status.Conditions) && r.DeletionTimestamp.IsZero() {
		if r.Spec.Networking.UserManagedNetworking == nil {
			userManagedNetworking := isNonePlatformOrSNO(r)
			//highAvailabilityMode := getHighAvailabilityMode(r, nil)
			// userManagedNetworking gets patched in one of two cases:
			// 1. Cluster topology is SNO.
			// 2. Platform is specified and platform is None, or External
			if !isNonePlatformOrSNO(r) && r.Spec.PlatformType != "" {
				platform := PlatformTypeToPlatform(r.Spec.PlatformType)
				/* platformUserManagedNetworking, err := webhooks.GetUMNFromPlatformAndHA(platform, highAvailabilityMode)
				if err != nil {
					agentclusterinstalllog.Info("failed to set UserManagedNetworking automatically due to error: %s", err.Error())
					return
				} */
				if platform != nil && (*platform.Type == models.PlatformTypeExternal || *platform.Type == models.PlatformTypeNone) {
					userManagedNetworking = true
				}
				//userManagedNetworking = platformUserManagedNetworking
			}
			r.Spec.Networking.UserManagedNetworking = &userManagedNetworking
			agentclusterinstalllog.Info("Setting UserManagedNetworking to %b", userManagedNetworking)
		}
	}
}

func PlatformTypeToPlatform(platformType PlatformType) *models.Platform {
	modelPlatformType := models.PlatformType(strings.ToLower(string(platformType)))
	platform := &models.Platform{Type: &modelPlatformType}
	return platform
}

//+kubebuilder:webhook:path=/validate-hiveextension-openshift-io-v1beta1-agentclusterinstall,mutating=false,failurePolicy=fail,sideEffects=None,groups=hiveextension.openshift.io,resources=agentclusterinstalls,verbs=create;update,versions=v1beta1,name=vagentclusterinstall.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AgentClusterInstall{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateCreate() (admission.Warnings, error) {
	agentclusterinstalllog.Info("validate create", "name", r.Name)
	// verify that UserNetworkManagement is not set to false with SNO.
	// if the user leave this field empty it is fine because the AI knows
	// what to set as default
	if isUserManagedNetworkingSetToFalseWithSNO(r) {
		err := fmt.Errorf("failed validation: UserManagedNetworking must be set to true with SNO")
		agentclusterinstalllog.Info(err.Error())
		return nil, err
	}
	/* 	if err := validateCreatePlatformAndUMN(r); err != nil {
		err = fmt.Errorf("Failed validation: %s", err.Error())
		agentclusterinstalllog.Info(err.Error())
		return nil, err

	} */
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	agentclusterinstalllog.Info("validate update", "name", r.Name)
	oldObject, ok := old.(*AgentClusterInstall)
	if !ok {
		return nil, fmt.Errorf("old object is not an AgentClusterInstall")
	}
	if !areImageSetRefsEqual(oldObject.Spec.ImageSetRef, r.Spec.ImageSetRef) {
		err := fmt.Errorf("Failed validation: Attempted to change AgentClusterInstall.ImageSetRef which is immutable")
		agentclusterinstalllog.Info(err.Error())
		return nil, err
	}

	if isUserManagedNetworkingSetToFalseWithSNO(r) {
		err := fmt.Errorf("Failed validation: UserManagedNetworking must be set to true with SNO")
		agentclusterinstalllog.Info(err.Error())
		return nil, err
	}

	/* 	if err := validateUpdatePlatformAndUMNUpdate(oldObject, r); err != nil {
		err := fmt.Errorf("Failed validation: ")
		agentclusterinstalllog.Info(err.Error())
		return nil, err
	} */

	if installAlreadyStarted(r.Status.Conditions) {
		ignoreChanges := mutableFields
		// MGMT-12794 This function returns true if the ProvisionRequirements field
		// has changed after installation completion. A change to this section has no effect
		// at this stage, but it is needed to serve some CI/CD gitops flows.
		if installCompleted(r.Status.Conditions) {
			ignoreChanges = append(ignoreChanges, "ProvisionRequirements")
		}
		hasChangedImmutableField, unsupportedDiff := hasChangedImmutableField(&oldObject.Spec, &r.Spec, ignoreChanges)
		if hasChangedImmutableField {
			err := fmt.Errorf("Failed validation: Attempted to change AgentClusterInstall.Spec which is immutable after install started, except for %s fields. Unsupported change: \n%s", strings.Join(mutableFields, ","), unsupportedDiff)
			agentclusterinstalllog.Info(err.Error())
			return nil, err
		}
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AgentClusterInstall) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func installAlreadyStarted(conditions []hivev1.ClusterInstallCondition) bool {
	cond := FindStatusCondition(conditions, ClusterCompletedCondition)
	if cond == nil {
		return false
	}
	switch cond.Reason {
	case ClusterInstallationFailedReason, ClusterInstalledReason, ClusterInstallationInProgressReason, ClusterAlreadyInstallingReason:
		return true
	default:
		return false
	}
}

func installCompleted(conditions []hivev1.ClusterInstallCondition) bool {
	cond := FindStatusCondition(conditions, ClusterCompletedCondition)
	if cond == nil {
		return false
	}
	return cond.Reason == ClusterInstalledReason || cond.Reason == ClusterInstallationFailedReason
}

// FindStatusCondition finds the conditionType in conditions.
func FindStatusCondition(conditions []hivev1.ClusterInstallCondition, conditionType hivev1.ClusterInstallConditionType) *hivev1.ClusterInstallCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

// hasChangedImmutableField determines if a AgentClusterInstall.spec immutable field was changed.
// it returns the diff string that shows the changes that are not supported
func hasChangedImmutableField(oldObject, cd *AgentClusterInstallSpec, mutableFields []string) (bool, string) {
	r := &diffReporter{}
	opts := cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.IgnoreFields(AgentClusterInstallSpec{}, mutableFields...),
		cmp.Reporter(r),
	}
	return !cmp.Equal(oldObject, cd, opts), r.String()
}

func areImageSetRefsEqual(imageSetRef1 *hivev1.ClusterImageSetReference, imageSetRef2 *hivev1.ClusterImageSetReference) bool {
	if imageSetRef1 == nil && imageSetRef2 == nil {
		return true
	} else if imageSetRef1 != nil && imageSetRef2 != nil {
		return imageSetRef1.Name == imageSetRef2.Name
	} else {
		return false
	}
}

func isUserManagedNetworkingSetToFalseWithSNO(newObject *AgentClusterInstall) bool {
	return isSNO(newObject) &&
		newObject.Spec.Networking.UserManagedNetworking != nil &&
		!*newObject.Spec.Networking.UserManagedNetworking
}

/* func validateCreatePlatformAndUMN(newObject *AgentClusterInstall) error {
	platform := PlatformTypeToPlatform(newObject.Spec.PlatformType)
	_, _, err := webhooks.GetActualCreateClusterPlatformParams(
		platform, newObject.Spec.Networking.UserManagedNetworking, getHighAvailabilityMode(newObject, nil), "")
	return err
} */

/* func validateUpdatePlatformAndUMNUpdate(oldObject, newObject *AgentClusterInstall) error {
	var (
		platform              *models.Platform
		userManagedNetworking *bool
	)

	if newObject.Spec.PlatformType != "" {
		platform = PlatformTypeToPlatform(newObject.Spec.PlatformType)
	} else {
		platform = PlatformTypeToPlatform(oldObject.Spec.PlatformType)
	}

	if newObject.Spec.Networking.UserManagedNetworking != nil {
		userManagedNetworking = newObject.Spec.Networking.UserManagedNetworking
	} else {
		userManagedNetworking = oldObject.Spec.Networking.UserManagedNetworking
	}

	_, _, err := webhooks.GetActualCreateClusterPlatformParams(
		platform, userManagedNetworking, getHighAvailabilityMode(oldObject, newObject), "")
	return err
} */

func getHighAvailabilityMode(originalObject, updatesObject *AgentClusterInstall) *string {
	if originalObject == nil {
		return swag.String("")
	}

	controlPlaneAgents := originalObject.Spec.ProvisionRequirements.ControlPlaneAgents
	workerAgents := originalObject.Spec.ProvisionRequirements.WorkerAgents

	if updatesObject != nil {
		if controlPlaneAgents != updatesObject.Spec.ProvisionRequirements.ControlPlaneAgents {
			controlPlaneAgents = updatesObject.Spec.ProvisionRequirements.ControlPlaneAgents
		}
		if workerAgents != updatesObject.Spec.ProvisionRequirements.WorkerAgents {
			workerAgents = updatesObject.Spec.ProvisionRequirements.WorkerAgents
		}
	}

	if controlPlaneAgents == 1 && workerAgents == 0 { // SNO
		return swag.String(models.ClusterHighAvailabilityModeNone)
	}
	return swag.String(models.ClusterHighAvailabilityModeFull)
}

func isSNO(newObject *AgentClusterInstall) bool {
	return newObject.Spec.ProvisionRequirements.ControlPlaneAgents == 1 &&
		newObject.Spec.ProvisionRequirements.WorkerAgents == 0
}

func isNonePlatformOrSNO(newObject *AgentClusterInstall) bool {
	return isSNO(newObject) && (newObject.Spec.PlatformType == "" || newObject.Spec.PlatformType == NonePlatformType) ||
		newObject.Spec.PlatformType == NonePlatformType
}

// diffReporter is a simple custom reporter that only records differences
// detected during comparison.
type diffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		p := r.path.String()
		vx, vy := r.path.Last().Values()
		r.diffs = append(r.diffs, fmt.Sprintf("\t%s: (%+v => %+v)", p, vx, vy))
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

/* func GetUMNFromPlatformAndHA(platform *models.Platform, highAvailabilityMode *string) (bool, error) {
	_, platformUserManagedNetworking, err := webhooks.GetClusterPlatformByHighAvailabilityMode(platform, nil, highAvailabilityMode)
	if err != nil {
		return false, err
	}
	return swag.BoolValue(platformUserManagedNetworking), nil
} */
