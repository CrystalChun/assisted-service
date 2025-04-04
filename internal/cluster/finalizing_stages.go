package cluster

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/openshift/assisted-service/models"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

const (
	generalWaitTimeout = 70 * time.Minute
	longWaitTimeout    = 10 * time.Hour
	shortWaitTimeout   = 10 * time.Minute
)

var finalizingStagesTimeoutsDefaultsHardTimeous = map[models.FinalizingStage]time.Duration{
	models.FinalizingStageWaitingForClusterOperators:              longWaitTimeout,
	models.FinalizingStageAddingRouterCa:                          generalWaitTimeout,
	models.FinalizingStageApplyingOlmManifests:                    shortWaitTimeout,
	models.FinalizingStageWaitingForOlmOperatorsCsv:               generalWaitTimeout,
	models.FinalizingStageWaitingForOlmOperatorsCsvInitialization: generalWaitTimeout,
	models.FinalizingStageWaitingForOLMOperatorSetupJobs:          shortWaitTimeout,
	models.FinalizingStageDone:                                    generalWaitTimeout,
}

var finalizingStagesTimeoutsDefaultsSoftTimeouts = map[models.FinalizingStage]time.Duration{
	models.FinalizingStageWaitingForClusterOperators:              generalWaitTimeout,
	models.FinalizingStageAddingRouterCa:                          generalWaitTimeout,
	models.FinalizingStageApplyingOlmManifests:                    shortWaitTimeout,
	models.FinalizingStageWaitingForOlmOperatorsCsv:               generalWaitTimeout,
	models.FinalizingStageWaitingForOlmOperatorsCsvInitialization: generalWaitTimeout,
	models.FinalizingStageWaitingForOLMOperatorSetupJobs:          shortWaitTimeout,
	models.FinalizingStageDone:                                    generalWaitTimeout,
}

var finalizingStages = []models.FinalizingStage{
	models.FinalizingStageWaitingForClusterOperators,
	models.FinalizingStageAddingRouterCa,
	models.FinalizingStageApplyingOlmManifests,
	models.FinalizingStageWaitingForOlmOperatorsCsvInitialization,
	models.FinalizingStageWaitingForOlmOperatorsCsv,
	models.FinalizingStageWaitingForOLMOperatorSetupJobs,
	models.FinalizingStageDone,
}

var nonFailingFinalizingStages = []models.FinalizingStage{
	models.FinalizingStageApplyingOlmManifests,
	models.FinalizingStageWaitingForOlmOperatorsCsvInitialization,
	models.FinalizingStageWaitingForOlmOperatorsCsv,
}

var olmOperatorFinalizingStages = []models.FinalizingStage{
	models.FinalizingStageWaitingForOlmOperatorsCsvInitialization,
	models.FinalizingStageWaitingForOlmOperatorsCsv,
	models.FinalizingStageWaitingForOLMOperatorSetupJobs,
}

func convertStageToEnvVar(stage models.FinalizingStage) string {
	return fmt.Sprintf("FINALIZING_STAGE_%s_TIMEOUT", strings.ReplaceAll(strings.ToUpper(string(stage)), " ", "_"))
}

func finalizingStageDefaultTimeout(stage models.FinalizingStage, softTimeoutEnabled bool, log logrus.FieldLogger) time.Duration {
	var (
		d   time.Duration
		err error
		ok  bool
	)
	val := os.Getenv(convertStageToEnvVar(stage))
	if val != "" {
		d, err = time.ParseDuration(val)
		if err == nil {
			return d
		}
		log.WithError(err).Warningf("failed to parse duration '%s' for stage '%s'", val, stage)
	}
	d, ok = lo.Ternary(softTimeoutEnabled, finalizingStagesTimeoutsDefaultsSoftTimeouts, finalizingStagesTimeoutsDefaultsHardTimeous)[stage]
	if ok {
		return d
	}
	log.Warningf("failed to get default timeout for stage '%s'", stage)
	return generalWaitTimeout
}

func finalizingStageTimeout(stage models.FinalizingStage, operators []*models.MonitoredOperator, softTimeoutEnabled bool, log logrus.FieldLogger) time.Duration {
	timeout := finalizingStageDefaultTimeout(stage, softTimeoutEnabled, log)
	if funk.Contains(olmOperatorFinalizingStages, stage) {
		timeoutSeconds := timeout.Seconds()
		for _, m := range operators {
			if m.OperatorType == models.OperatorTypeOlm {
				timeoutSeconds = math.Max(timeoutSeconds, float64(m.TimeoutSeconds))
			}
		}
		timeout = time.Duration(timeoutSeconds) * time.Second
	}
	return timeout
}
