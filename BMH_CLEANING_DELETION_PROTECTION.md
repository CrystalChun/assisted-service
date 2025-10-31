# BMH Automated Cleaning Deletion Protection

## Overview

This modification ensures that when a BareMetalHost (BMH) has `automatedCleaningMode` set to `metadata` or any value other than `disabled`, the associated InfraEnv and PreprovisioningImage resources cannot be deleted until the BMH is removed or has completed deprovisioning.

## Why This Is Needed

When `automatedCleaningMode: metadata` is set on a BMH:
1. During deprovisioning, Ironic boots the Ironic Python Agent (IPA) from the PreprovisioningImage
2. IPA performs metadata cleaning (erases partition tables, filesystem signatures, etc.)
3. If the PreprovisioningImage or InfraEnv is deleted before deprovisioning completes, the BMH will fail to clean properly

## Changes Made

### 1. InfraEnv Controller (`infraenv_controller.go`)

#### Added Import
```go
bmh_v1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
```

#### Added Constant
```go
bmhInfraEnvLabel = "infraenvs.agent-install.openshift.io"
```

#### Modified `deregisterInfraEnvWithHosts()` Method
- Added check for BMHs requiring cleaning resources before allowing InfraEnv deletion
- Calls new `checkBMHsRequiringCleaning()` helper function

#### Added `checkBMHsRequiringCleaning()` Method
- Lists all BareMetalHosts in the InfraEnv's namespace
- Filters for BMHs that reference the InfraEnv via the `infraenvs.agent-install.openshift.io` label
- Checks if BMH has `automatedCleaningMode != disabled`
- Checks if BMH is in a state that requires the image:
  - `StateProvisioned`
  - `StateProvisioning`
  - `StateDeprovisioning`
  - `StateDeleting`
  - `StatePreparing`
  - `StateReady`
  - `StateAvailable`
  - `StateInspecting`
- Returns an error with details if any BMHs still require the InfraEnv

### 2. PreprovisioningImage Controller (`preprovisioningimage_controller.go`)

#### Added Import
```go
"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
```

#### Added Constants
```go
PreprovisioningImageFinalizerName = "preprovisioningimage." + aiv1beta1.Group + "/ai-deprovision"
preprovisioningImageInfraEnvLabel = "infraenvs.agent-install.openshift.io"
```

#### Added RBAC Permissions
```go
// +kubebuilder:rbac:groups=metal3.io,resources=baremetalhosts,verbs=get;list;watch
```

#### Modified `Reconcile()` Method
- Added deletion timestamp check at the beginning
- Added finalizer management
- Calls `handlePreprovisioningImageDeletion()` when resource is being deleted
- Calls `ensurePreprovisioningImageFinalizer()` to add finalizer

#### Added `ensurePreprovisioningImageFinalizer()` Method
- Adds finalizer to PreprovisioningImage if not already present

#### Added `handlePreprovisioningImageDeletion()` Method
- Gets the BMH that owns the PreprovisioningImage
- Checks if BMH has `automatedCleaningMode != disabled`
- Checks if BMH is in a state that requires the image (same states as InfraEnv check)
- Returns an error if BMH still requires the image
- Removes finalizer when safe to delete

## Error Messages

### InfraEnv Deletion Blocked
```
cannot delete InfraEnv <namespace>/<name>: <N> BareMetalHost(s) with automatedCleaningMode enabled still reference it: [<bmh-namespace>/<bmh-name> (state: <state>, cleaning: <mode>), ...]. These hosts require the PreprovisioningImage for deprovisioning. Please remove or deprovision these hosts first.
```

### PreprovisioningImage Deletion Blocked
```
cannot delete PreprovisioningImage <namespace>/<name>: BareMetalHost <bmh-namespace>/<bmh-name> with automatedCleaningMode=<mode> (state: <state>) still requires it for deprovisioning. Please remove or deprovision the host first.
```

## How to Test

### Test Case 1: Normal Deletion Flow
1. Create an InfraEnv
2. Create a BMH with `automatedCleaningMode: metadata` referencing the InfraEnv
3. Provision the BMH
4. Delete the BMH
5. Wait for BMH to be fully deleted
6. Delete the InfraEnv - should succeed

### Test Case 2: Blocked Deletion
1. Create an InfraEnv
2. Create a BMH with `automatedCleaningMode: metadata` referencing the InfraEnv
3. Provision the BMH
4. Try to delete the InfraEnv while BMH is still provisioned
5. Deletion should be blocked with error message
6. Delete the BMH first
7. Then delete the InfraEnv - should succeed

### Test Case 3: Disabled Cleaning
1. Create an InfraEnv
2. Create a BMH with `automatedCleaningMode: disabled` referencing the InfraEnv
3. Provision the BMH
4. Try to delete the InfraEnv
5. Deletion should succeed (cleaning disabled, no IPA needed)

## Backward Compatibility

- Existing InfraEnv and PreprovisioningImage resources will get finalizers added automatically
- No breaking changes to existing APIs
- Only affects deletion behavior when BMHs with cleaning enabled are present

## Future Enhancements

Consider adding:
- Status conditions on InfraEnv/PreprovisioningImage indicating deletion is blocked
- Events to notify users why deletion is blocked
- Webhook validation to provide immediate feedback on attempted deletion

