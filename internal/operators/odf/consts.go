package odf

import "github.com/openshift/assisted-service/models"

type odfDeploymentMode string

const (
	compactMode  odfDeploymentMode = "Compact"  // only masters, control plane nodes will run ODF
	standardMode odfDeploymentMode = "Standard" // at least 3 masters and 3 workers, workers will run ODF
	unknown      odfDeploymentMode = "Unknown"  // none of the above, the mode is not determined yet

	clusterValidationID = string(models.ClusterValidationIDOdfRequirementsSatisfied)
	hostValidationID    = string(models.HostValidationIDOdfRequirementsSatisfied)
)
