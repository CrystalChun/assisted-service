package featuresupport

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
	"github.com/sirupsen/logrus"
)

func getPlatformFilters() []SupportLevelFilters {
	return []SupportLevelFilters{
		{PlatformType: models.PlatformTypeVsphere.Pointer()},
		{PlatformType: models.PlatformTypeNutanix.Pointer()},
		{PlatformType: models.PlatformTypeBaremetal.Pointer()},
		{PlatformType: models.PlatformTypeNone.Pointer()},
		{PlatformType: models.PlatformTypeExternal.Pointer()},
		{
			PlatformType:         models.PlatformTypeExternal.Pointer(),
			ExternalPlatformName: swag.String(common.ExternalPlatformNameOci),
		},
	}
}

var _ = Describe("V2ListFeatureSupportLevels API", func() {
	availableVersions := []string{"4.9", "4.10", "4.11", "4.12", "4.13"}
	availableCpuArch := []string{
		models.ClusterCPUArchitectureX8664,
		models.ClusterCPUArchitectureArm64,
		models.ClusterCreateParamsCPUArchitectureAarch64,
		models.ClusterCPUArchitectureS390x,
		models.ClusterCPUArchitecturePpc64le,
		models.ClusterCPUArchitectureMulti,
	}

	Context("Feature compatibility", func() {
		for _, f := range featuresList {
			for _, v := range availableVersions {
				for _, a := range availableCpuArch {
					feature := f
					version := v
					arch := a

					It(fmt.Sprintf("isFeatureCompatibleWithArchitecture %s, %s, %s", version, feature, arch), func() {
						filters := SupportLevelFilters{OpenshiftVersion: version, CPUArchitecture: swag.String(arch)}
						isFeatureCompatibleWithArchitecture(feature, filters.OpenshiftVersion, swag.StringValue(filters.CPUArchitecture))
					})
				}
			}
		}
	})

	It("Test ARM64 is not supported under 4.10", func() {
		feature := models.ArchitectureSupportLevelIDARM64ARCHITECTURE
		Expect(isArchitectureSupported(feature, "4.6")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.7")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.8")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.9")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.10")).To(BeTrue())
		Expect(isArchitectureSupported(feature, "4.11")).To(BeTrue())
		Expect(isArchitectureSupported(feature, "4.12")).To(BeTrue())
		Expect(isArchitectureSupported(feature, "4.13")).To(BeTrue())

		// Check for feature release
		Expect(isArchitectureSupported(feature, "4.30")).To(BeTrue())
	})

	It("Test s390x is not supported under 4.12", func() {
		feature := models.ArchitectureSupportLevelIDS390XARCHITECTURE
		Expect(isArchitectureSupported(feature, "4.6")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.7")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.8")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.9")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.10")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.11")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.12")).To(BeTrue())
		Expect(isArchitectureSupported(feature, "4.13")).To(BeTrue())

		// Check for feature release
		Expect(isArchitectureSupported(feature, "4.30")).To(BeTrue())

	})

	It("Test PPC64LE is not supported under 4.12", func() {
		feature := models.ArchitectureSupportLevelIDPPC64LEARCHITECTURE
		Expect(isArchitectureSupported(feature, "4.6")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.7")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.8")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.9")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.10")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.11")).To(BeFalse())
		Expect(isArchitectureSupported(feature, "4.12")).To(BeTrue())
		Expect(isArchitectureSupported(feature, "4.13")).To(BeTrue())

		// Check for feature release
		Expect(isArchitectureSupported(feature, "4.30")).To(BeTrue())
	})

	Context("Test LSO CPU compatibility", func() {
		feature := models.FeatureSupportLevelIDLSO
		It("LSO IsFeatureAvailable", func() {
			Expect(IsFeatureAvailable(feature, "Does not matter", swag.String(models.ClusterCPUArchitecturePpc64le))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "Does not matter", swag.String(models.ClusterCPUArchitectureX8664))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "Does not matter", swag.String(models.ClusterCPUArchitectureS390x))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "Does not matter", swag.String(models.ClusterCPUArchitectureArm64))).To(BeFalse())
		})
		It("LSO GetSupportLevel on architecture", func() {
			featureSupportParams := SupportLevelFilters{OpenshiftVersion: "Any", CPUArchitecture: swag.String(models.ClusterCPUArchitectureX8664)}
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelSupported))

			featureSupportParams.CPUArchitecture = swag.String(models.ClusterCPUArchitectureS390x)
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelSupported))

			featureSupportParams.CPUArchitecture = swag.String(models.ClusterCPUArchitecturePpc64le)
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelSupported))

			featureSupportParams.CPUArchitecture = swag.String(models.ClusterCPUArchitectureArm64)
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelUnavailable))
		})
	})

	Context("Test non-standard HA OCP Control Plane", func() {
		feature := models.FeatureSupportLevelIDNONSTANDARDHACONTROLPLANE
		arch := common.X86CPUArchitecture

		It("test feature availability", func() {
			Expect(IsFeatureAvailable(feature, common.MinimumVersionForNonStandardHAOCPControlPlane, swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.17", swag.String(arch))).To(BeFalse())
		})

		DescribeTable("test feature compatability with other features", func(activeFeatures []SupportLevelFeature, shouldSucceed bool) {
			activeFeatures = append(activeFeatures, &NonStandardHAControlPlane{})

			if shouldSucceed {
				Expect(
					isFeaturesCompatibleWithFeatures(
						common.MinimumVersionForNonStandardHAOCPControlPlane,
						activeFeatures),
				).ToNot(HaveOccurred())
			} else {
				Expect(
					isFeaturesCompatibleWithFeatures(
						common.MinimumVersionForNonStandardHAOCPControlPlane,
						activeFeatures),
				).To(HaveOccurred())
			}
		},
			Entry(
				"platform baremetal",
				[]SupportLevelFeature{&BaremetalPlatformFeature{}},
				true,
			),

			Entry(
				"none platform",
				[]SupportLevelFeature{&NonePlatformFeature{}},
				true,
			),

			Entry(
				"external platform",
				[]SupportLevelFeature{&ExternalPlatformFeature{}},
				false,
			),

			Entry(
				"nutanix platform",
				[]SupportLevelFeature{&NutanixIntegrationFeature{}},
				false,
			),

			Entry(
				"vsphere platform",
				[]SupportLevelFeature{&VsphereIntegrationFeature{}},
				false,
			),
		)

		DescribeTable(
			"test feature architecture support",
			func(arch string, result bool) {
				Expect(
					isFeatureCompatibleWithArchitecture(
						&NonStandardHAControlPlane{},
						common.MinimumVersionForNonStandardHAOCPControlPlane,
						arch,
					),
				).To(Equal(result))
			},
			Entry(
				common.X86CPUArchitecture,
				common.X86CPUArchitecture,
				true,
			),

			Entry(
				common.ARM64CPUArchitecture,
				common.ARM64CPUArchitecture,
				false,
			),

			Entry(
				common.S390xCPUArchitecture,
				common.S390xCPUArchitecture,
				false,
			),

			Entry(
				common.PowerCPUArchitecture,
				common.PowerCPUArchitecture,
				false,
			),

			Entry(
				common.MultiCPUArchitecture,
				common.MultiCPUArchitecture,
				false,
			),
		)
	})

	Context("Test TNA", func() {
		feature := models.FeatureSupportLevelIDTNA
		openshiftVersionWithoutTNA := "4.18"
		openshiftVersionSupportWithTNA := "4.20"

		It("test feature availability", func() {
			Expect(IsFeatureAvailable(feature, common.MinimumVersionForArbiterClusters, nil)).To(BeTrue())
			Expect(IsFeatureAvailable(feature, openshiftVersionWithoutTNA, nil)).To(BeFalse())
		})

		DescribeTable("test feature compatability with other features", func(activeFeatures []SupportLevelFeature, shouldSucceed bool) {
			activeFeatures = append(activeFeatures, &TnaFeature{})

			if shouldSucceed {
				Expect(
					isFeaturesCompatibleWithFeatures(
						common.MinimumVersionForArbiterClusters,
						activeFeatures),
				).ToNot(HaveOccurred())
			} else {
				Expect(
					isFeaturesCompatibleWithFeatures(
						common.MinimumVersionForArbiterClusters,
						activeFeatures),
				).To(HaveOccurred())
			}
		},
			Entry(
				"platform baremetal",
				[]SupportLevelFeature{&BaremetalPlatformFeature{}},
				true,
			),

			Entry(
				"none platform",
				[]SupportLevelFeature{&NonePlatformFeature{}},
				false,
			),

			Entry(
				"external platform",
				[]SupportLevelFeature{&ExternalPlatformFeature{}},
				false,
			),

			Entry(
				"nutanix platform",
				[]SupportLevelFeature{&NutanixIntegrationFeature{}},
				false,
			),

			Entry(
				"vsphere platform",
				[]SupportLevelFeature{&VsphereIntegrationFeature{}},
				false,
			),
		)

		DescribeTable(
			"test feature architecture support",
			func(arch string, result bool) {
				Expect(
					isFeatureCompatibleWithArchitecture(
						&TnaFeature{},
						common.MinimumVersionForArbiterClusters,
						arch,
					),
				).To(Equal(result))
			},
			Entry(
				common.X86CPUArchitecture,
				common.X86CPUArchitecture,
				true,
			),

			Entry(
				common.ARM64CPUArchitecture,
				common.ARM64CPUArchitecture,
				true,
			),

			Entry(
				common.S390xCPUArchitecture,
				common.S390xCPUArchitecture,
				true,
			),

			Entry(
				common.PowerCPUArchitecture,
				common.PowerCPUArchitecture,
				true,
			),

			Entry(
				common.MultiCPUArchitecture,
				common.MultiCPUArchitecture,
				true,
			),
		)

		DescribeTable(
			"test feature active level",
			func(cluster *common.Cluster, result featureActiveLevel) {
				Expect(
					(&TnaFeature{}).getFeatureActiveLevel(cluster, nil, nil, nil),
				).To(Equal(result))
			},
			Entry(
				"active",
				&common.Cluster{
					Cluster: models.Cluster{
						OpenshiftVersion: common.MinimumVersionForArbiterClusters,
						Hosts: []*models.Host{
							{Role: models.HostRoleArbiter},
						},
						Platform: &models.Platform{
							Type: models.PlatformTypeBaremetal.Pointer(),
						},
					},
				},
				activeLevelActive,
			),

			Entry(
				"not active - nil cluster",
				nil,
				activeLevelNotActive,
			),

			Entry(
				"not active - cluster doesn't have arbiter nodes",
				&common.Cluster{
					Cluster: models.Cluster{
						OpenshiftVersion: common.MinimumVersionForArbiterClusters,
						Hosts: []*models.Host{
							{Role: models.HostRoleMaster},
						},
						Platform: &models.Platform{
							Type: models.PlatformTypeBaremetal.Pointer(),
						},
					},
				},
				activeLevelNotActive,
			),
		)

		DescribeTable(
			"test feature support level",
			func(filters SupportLevelFilters, result models.SupportLevel) {
				supportLevel, _ := (&TnaFeature{}).getSupportLevel(filters)
				Expect(supportLevel).To(Equal(result))
			},
			Entry(
				"tech preview openshift version with platform filter",
				SupportLevelFilters{
					OpenshiftVersion: common.MinimumVersionForArbiterClusters,
					PlatformType:     models.PlatformTypeBaremetal.Pointer(),
				},
				models.SupportLevelTechPreview,
			),

			Entry(
				"tech preview openshift version without platform filter",
				SupportLevelFilters{
					OpenshiftVersion: common.MinimumVersionForArbiterClusters,
				},
				models.SupportLevelTechPreview,
			),

			Entry(
				"unavailable - platform is not baremetal",
				SupportLevelFilters{
					OpenshiftVersion: common.MinimumVersionForArbiterClusters,
					PlatformType:     models.PlatformTypeNone.Pointer(),
				},
				models.SupportLevelUnavailable,
			),

			Entry(
				"unavailable - openshift version is too low",
				SupportLevelFilters{
					OpenshiftVersion: openshiftVersionWithoutTNA,
				},
				models.SupportLevelUnavailable,
			),

			Entry(
				"support openshift version with platform filter",
				SupportLevelFilters{
					OpenshiftVersion: openshiftVersionSupportWithTNA,
					PlatformType:     models.PlatformTypeBaremetal.Pointer(),
				},
				models.SupportLevelSupported,
			),

			Entry(
				"support openshift version without platform filter",
				SupportLevelFilters{
					OpenshiftVersion: openshiftVersionSupportWithTNA,
				},
				models.SupportLevelSupported,
			),
		)
	})

	Context("Test MCE not supported under 4.10", func() {
		feature := models.FeatureSupportLevelIDMCE
		It(fmt.Sprintf("%s test", feature), func() {
			arch := "DoesNotMatter"
			Expect(IsFeatureAvailable(feature, "4.9", swag.String(arch))).To(BeFalse())
			Expect(IsFeatureAvailable(feature, "4.10", swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.11", swag.String(arch))).To(BeTrue())

			featureSupportParams := SupportLevelFilters{OpenshiftVersion: "4.9", CPUArchitecture: swag.String(arch)}
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelUnavailable))
			featureSupportParams = SupportLevelFilters{OpenshiftVersion: "4.11.20", CPUArchitecture: swag.String(arch)}
			Expect(GetSupportLevel(feature, featureSupportParams)).To(Equal(models.SupportLevelSupported))

			Expect(IsFeatureAvailable(feature, "4.12", swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.13", swag.String(arch))).To(BeTrue())

			// Check for feature release
			Expect(IsFeatureAvailable(feature, "4.30", swag.String(arch))).To(BeTrue())
		})
	})

	Context("Test network type", func() {
		It("Test SDN not supported over 4.15", func() {
			feature := models.FeatureSupportLevelIDSDNNETWORKTYPE
			arch := "DoesNotMatter"
			Expect(IsFeatureAvailable(feature, "4.14", swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.15", swag.String(arch))).To(BeFalse())
			Expect(IsFeatureAvailable(feature, "4.16", swag.String(arch))).To(BeFalse())
		})

		It("Test OVN is supported over 4.15", func() {
			feature := models.FeatureSupportLevelIDOVNNETWORKTYPE
			arch := "DoesNotMatter"
			Expect(IsFeatureAvailable(feature, "4.14", swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.15", swag.String(arch))).To(BeTrue())
			Expect(IsFeatureAvailable(feature, "4.16", swag.String(arch))).To(BeTrue())
		})
	})

	Context("GetCpuArchitectureSupportList", func() {
		It("GetCpuArchitectureSupportList for openshift version 4.6", func() {
			openshiftVersion := "4.6"
			supportedArchitectures := GetCpuArchitectureSupportList(openshiftVersion)
			Expect(len(supportedArchitectures)).To(Equal(5))

			for key, value := range supportedArchitectures {
				if key == string(models.ArchitectureSupportLevelIDX8664ARCHITECTURE) {
					Expect(value).To(Equal(models.SupportLevelSupported))
				} else {
					Expect(value).To(Equal(models.SupportLevelUnavailable))
				}
			}
		})

		It("GetCpuArchitectureSupportList for openshift version 4.13", func() {
			openshiftVersion := "4.13"
			supportedArchitectures := GetCpuArchitectureSupportList(openshiftVersion)
			Expect(len(supportedArchitectures)).To(Equal(5))
			for key, value := range supportedArchitectures {
				if key == string(models.ArchitectureSupportLevelIDMULTIARCHRELEASEIMAGE) {
					Expect(value).To(Equal(models.SupportLevelTechPreview))
				} else {
					Expect(value).To(Equal(models.SupportLevelSupported))
				}
			}
		})
	})

	Context("Test OCI platform support", func() {
		DescribeTable(
			"Validation pass",
			func(openshiftVersion string, expectedSupportLevel models.SupportLevel) {
				filters := SupportLevelFilters{
					OpenshiftVersion: openshiftVersion,
					CPUArchitecture:  swag.String(common.DefaultCPUArchitecture),
				}
				supportLevel := GetSupportLevel(models.FeatureSupportLevelIDEXTERNALPLATFORMOCI, filters)
				Expect(supportLevel).To(Equal(expectedSupportLevel))
			},
			Entry("OCI unavailable with Openshift 4.13", "4.13", models.SupportLevelUnavailable),
			Entry("OCI tech-preview with Openshift 4.14", "4.14", models.SupportLevelSupported),
			Entry("OCI tech-preview with Openshidt 4.15", "4.15", models.SupportLevelSupported),
		)
	})

	Context("GetSupportList", func() {

		for _, filters := range getPlatformFilters() {
			filters := filters
			When("GetFeatureSupportList 4.12 with Platform", func() {
				It(string(*filters.PlatformType)+" "+swag.StringValue(filters.ExternalPlatformName), func() {
					list := GetFeatureSupportList("dummy", nil, filters.PlatformType, filters.ExternalPlatformName)
					Expect(len(list)).To(Equal(41))
				})
			})
		}

		It("GetFeatureSupportList 4.12", func() {
			list := GetFeatureSupportList("4.12", nil, nil, nil)
			Expect(len(list)).To(Equal(46))
		})

		It("GetFeatureSupportList 4.13", func() {
			list := GetFeatureSupportList("4.13", nil, nil, nil)
			Expect(len(list)).To(Equal(46))
		})

		It("GetCpuArchitectureSupportList 4.12", func() {
			list := GetCpuArchitectureSupportList("4.12")
			Expect(len(list)).To(Equal(5))
		})

		It("GetCpuArchitectureSupportList 4.13", func() {
			list := GetCpuArchitectureSupportList("4.13")
			Expect(len(list)).To(Equal(5))
		})

		It("GetFeatureSupportList 4.11 with not supported architecture", func() {
			featuresList := GetFeatureSupportList("4.11", swag.String(models.ClusterCPUArchitecturePpc64le), nil, nil)

			for _, feature := range featuresList {
				Expect(feature.SupportLevel).To(Equal(models.SupportLevelUnavailable))
			}
		})

		It("GetFeatureSupportList 4.13 with unsupported architecture", func() {
			featuresList := GetFeatureSupportList("4.12", swag.String(models.ClusterCPUArchitecturePpc64le), nil, nil)
			for _, feature := range featuresList {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDSNO {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelUnavailable))
			}

			featuresList = GetFeatureSupportList("4.13", swag.String(models.ClusterCPUArchitecturePpc64le), nil, nil)
			for _, feature := range featuresList {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDSNO {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelDevPreview))
			}
		})

		It("GetFeatureSupportList 4.13 with supported architecture", func() {
			featuresList := GetFeatureSupportList("4.13", swag.String(models.ClusterCPUArchitectureX8664), nil, nil)
			for _, feature := range featuresList {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDSNO {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelSupported))
			}
		})
	})

	Context("ValidateIncompatibleFeatures", func() {
		log := logrus.New()

		It("No feature is activated", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion: "4.6",
				CPUArchitecture:  models.ClusterCPUArchitectureX8664,
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureX8664, &cluster, nil, nil)).To(BeNil())
		})

		It("No OCP version with CPU architecture that depends on OCP version", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				CPUArchitecture:       models.ClusterCPUArchitectureArm64,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureArm64, &cluster, nil, nil)).To(BeNil())
		})
		It("Single compatible feature is activated", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.8",
				CPUArchitecture:       models.ClusterCPUArchitectureX8664,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureX8664, &cluster, nil, nil)).To(BeNil())
		})
		It("Update s390x cluster", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.8",
				CPUArchitecture:       models.ClusterCPUArchitectureS390x,
				ControlPlaneCount:     common.MinMasterHostsNeededForInstallationInHaMode,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			params := models.V2ClusterUpdateParams{UserManagedNetworking: swag.Bool(false)}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, nil, &params)).To(Not(BeNil()))
		})
		It("Ignore validation on AddHostCluster", func() {
			logBuffer := bytes.Buffer{}
			testLogger := logrus.New()
			testLogger.SetOutput(&logBuffer)

			clusterID := strfmt.UUID("e679ea3f-3b85-40e0-8dc9-82fd6945d9b2")
			cluster := common.Cluster{Cluster: models.Cluster{
				ID:                &clusterID,
				OpenshiftVersion:  "4.19",
				Kind:              swag.String(models.ClusterKindAddHostsCluster),
				CPUArchitecture:   models.ClusterCPUArchitectureS390x,
				ControlPlaneCount: common.MinMasterHostsNeededForInstallationInHaMode,
				Platform:          &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}

			infraEnv := models.InfraEnv{
				ClusterID:       *cluster.ID,
				CPUArchitecture: models.ClusterCPUArchitectureS390x,
			}

			Expect(ValidateActiveFeatures(logrus.NewEntry(testLogger), &cluster, &infraEnv, nil)).To(BeNil())
			Expect(logBuffer.String()).To(ContainSubstring("skipping feature support validation:"))
		})
		It("Update s390x cluster", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.13",
				CPUArchitecture:       models.ClusterCPUArchitectureS390x,
				ControlPlaneCount:     common.MinMasterHostsNeededForInstallationInHaMode,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			infraEnv := models.InfraEnv{Type: common.ImageTypePtr(models.ImageTypeFullIso)}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, &infraEnv, nil)).To(BeNil())

			params := models.InfraEnvUpdateParams{ImageType: models.ImageTypeMinimalIso}
			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, &infraEnv, &params)
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(ContainSubstring("cannot use Minimal ISO because it's not compatible with the s390x architecture on version 4.13 of OpenShift"))
		})
		It("SNO feature is activated with incompatible architecture ppc64le on 4.12", func() {
			expectedError := "cannot use Single Node OpenShift because it's not compatible with the ppc64le architecture on version 4.12 of OpenShift"
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.12",
				CPUArchitecture:       models.ClusterCPUArchitecturePpc64le,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitecturePpc64le, &cluster, nil, nil).Error()).To(Equal(expectedError))
		})
		It("SNO feature is compatible on ppc64le architecture at 4.13", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.13",
				CPUArchitecture:       models.ClusterCPUArchitecturePpc64le,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitecturePpc64le, &cluster, nil, nil)).To(BeNil())
		})
		It("SNO feature is activated with incompatible architecture s390x on 4.12", func() {
			expectedError := "cannot use Single Node OpenShift because it's not compatible with the s390x architecture on version 4.12 of OpenShift"
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.12",
				CPUArchitecture:       models.ClusterCPUArchitectureS390x,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, nil, nil).Error()).To(Equal(expectedError))
		})
		It("SNO feature is activated with compatible architecture s390x on 4.13", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.13",
				CPUArchitecture:       models.ClusterCPUArchitectureS390x,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, nil, nil)).To(BeNil())
		})
		It("Nutanix feature is activated with incompatible architecture", func() {
			expectedError := "cannot use arm64 architecture because it's not compatible on version 4.8 of OpenShift"
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.8",
				CPUArchitecture:       models.ClusterCPUArchitectureArm64,
				ControlPlaneCount:     1,
				UserManagedNetworking: swag.Bool(true),
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureArm64, &cluster, nil, nil).Error()).To(Equal(expectedError))
		})
		It("ClusterManagedNetworking feature is activated with compatible architecture on 4.11", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:      "4.11",
				CPUArchitecture:       models.ClusterCPUArchitectureArm64,
				ControlPlaneCount:     common.MinMasterHostsNeededForInstallationInHaMode,
				Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeBaremetal)},
				UserManagedNetworking: swag.Bool(false),
			}}
			Expect(ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureArm64, &cluster, nil, nil)).To(BeNil())
		})
		It("Ppc64le with CMN - fail", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion: "4.12",
				CPUArchitecture:  models.ClusterCPUArchitecturePpc64le,
			}}
			infraEnv := models.InfraEnv{CPUArchitecture: models.ClusterCPUArchitecturePpc64le, Type: common.ImageTypePtr(models.ImageTypeFullIso)}

			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitecturePpc64le, &cluster, nil, nil)
			Expect(err).To(Not(BeNil()))
			cluster.UserManagedNetworking = swag.Bool(true)
			err = ValidateIncompatibleFeatures(log, models.ClusterCPUArchitecturePpc64le, &cluster, &infraEnv, nil)
			Expect(err).To(BeNil())
		})
		It("s390x with CMN and minimal iso - fail", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion: "4.12",
				CPUArchitecture:  models.ClusterCPUArchitectureS390x,
			}}
			infraEnv := models.InfraEnv{CPUArchitecture: models.ClusterCPUArchitectureS390x, Type: common.ImageTypePtr(models.ImageTypeMinimalIso)}

			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, nil, nil)
			Expect(err).To(Not(BeNil()))
			cluster.UserManagedNetworking = swag.Bool(true)
			err = ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, &infraEnv, nil)
			Expect(err).To(Not(BeNil()))
		})
		It("s390x with External and platformName=oci - fail", func() {
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion: "4.14",
				CPUArchitecture:  models.ClusterCPUArchitectureS390x,
				Platform: &models.Platform{
					Type: common.PlatformTypePtr(models.PlatformTypeExternal),
					External: &models.PlatformExternal{
						PlatformName:           swag.String("oci"),
						CloudControllerManager: swag.String(models.PlatformExternalCloudControllerManagerExternal),
					},
				},
				UserManagedNetworking: swag.Bool(true),
			}}
			infraEnv := models.InfraEnv{CPUArchitecture: models.ClusterCPUArchitectureS390x, Type: common.ImageTypePtr(models.ImageTypeFullIso)}

			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureS390x, &cluster, &infraEnv, nil)
			Expect(err).To(Not(BeNil()))
		})
		It("Nutanix with incompatible features - fail", func() {
			operatorsCNV := []*models.MonitoredOperator{
				{
					Name:             "cnv",
					Namespace:        "openshift-cnv",
					OperatorType:     models.OperatorTypeOlm,
					SubscriptionName: "hco-operatorhub",
					TimeoutSeconds:   60 * 60,
				},
			}
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:   "4.14",
				CPUArchitecture:    models.ClusterCPUArchitectureX8664,
				Platform:           &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
				MonitoredOperators: operatorsCNV,
			}}
			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureX8664, &cluster, nil, nil)
			Expect(err).To(HaveOccurred())

			operatorsMCE := []*models.MonitoredOperator{
				{
					Name:             "mce",
					OperatorType:     models.OperatorTypeOlm,
					Namespace:        "multicluster-engine",
					SubscriptionName: "multicluster-engine",
					TimeoutSeconds:   60 * 60,
				},
			}
			cluster = common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:   "4.14",
				CPUArchitecture:    models.ClusterCPUArchitectureX8664,
				Platform:           &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
				MonitoredOperators: operatorsMCE,
			}}
			err = ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureX8664, &cluster, nil, nil)
			Expect(err).To(HaveOccurred())
		})
		It("VSphere with incompatible features - fail", func() {
			operatorsCNV := []*models.MonitoredOperator{
				{
					Name:             "cnv",
					Namespace:        "openshift-cnv",
					OperatorType:     models.OperatorTypeOlm,
					SubscriptionName: "hco-operatorhub",
					TimeoutSeconds:   60 * 60,
				},
			}
			cluster := common.Cluster{Cluster: models.Cluster{
				OpenshiftVersion:   "4.14",
				CPUArchitecture:    models.ClusterCPUArchitectureX8664,
				Platform:           &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeVsphere)},
				MonitoredOperators: operatorsCNV,
			}}
			err := ValidateIncompatibleFeatures(log, models.ClusterCPUArchitectureX8664, &cluster, nil, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Incompatibilities", func() {
		Context("IsFeatureActivated", func() {
			It("Activated features in cluster - Sno, VipAutoAlloc, UserManagedNetworking, NutanixIntegration", func() {
				operators := []*models.MonitoredOperator{
					{
						Name:             "cnv",
						Namespace:        "openshift-cnv",
						OperatorType:     models.OperatorTypeOlm,
						SubscriptionName: "hco-operatorhub",
						TimeoutSeconds:   60 * 60,
					},
				}

				cluster := common.Cluster{Cluster: models.Cluster{
					OpenshiftVersion:      "4.8",
					CPUArchitecture:       models.ClusterCPUArchitecturePpc64le,
					ControlPlaneCount:     1,
					UserManagedNetworking: swag.Bool(true),
					Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
					VipDhcpAllocation:     swag.Bool(true),
					MonitoredOperators:    operators,
				},
				}

				activatedFeatures := []SupportLevelFeature{
					&VipAutoAllocFeature{}, &SnoFeature{}, &UserManagedNetworkingFeature{}, &NutanixIntegrationFeature{}, &CnvFeature{},
				}

				for _, feature := range activatedFeatures {
					Expect(feature.getFeatureActiveLevel(&cluster, nil, nil, nil)).To(Equal(activeLevelActive))
				}
			})

			It("Disable activated features in cluster - Sno, VipAutoAlloc, UserManagedNetworking, NutanixIntegration, Cnv", func() {
				operators := []*models.MonitoredOperator{
					{
						Name:             "cnv",
						Namespace:        "openshift-cnv",
						OperatorType:     models.OperatorTypeOlm,
						SubscriptionName: "hco-operatorhub",
						TimeoutSeconds:   60 * 60,
					},
				}

				cluster := common.Cluster{Cluster: models.Cluster{
					OpenshiftVersion:      "4.8",
					CPUArchitecture:       models.ClusterCPUArchitecturePpc64le,
					ControlPlaneCount:     1,
					UserManagedNetworking: swag.Bool(true),
					Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
					VipDhcpAllocation:     swag.Bool(true),
					MonitoredOperators:    operators,
				}}
				params := models.V2ClusterUpdateParams{
					VipDhcpAllocation:     swag.Bool(false),
					Platform:              &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeBaremetal)},
					UserManagedNetworking: swag.Bool(false),
					OlmOperators: []*models.OperatorCreateParams{
						{Name: "lvm"},
					},
				}

				activatedFeatures := []SupportLevelFeature{
					&VipAutoAllocFeature{}, &UserManagedNetworkingFeature{}, &NutanixIntegrationFeature{}, &CnvFeature{},
				}

				for _, feature := range activatedFeatures {
					Expect(feature.getFeatureActiveLevel(&cluster, nil, nil, nil)).To(Equal(activeLevelActive))
					Expect(feature.getFeatureActiveLevel(&cluster, nil, &params, nil)).To(Equal(activeLevelNotActive))
				}
				Expect((&SnoFeature{}).getFeatureActiveLevel(&cluster, nil, &params, nil)).To(Equal(activeLevelActive))
				Expect((&LvmFeature{}).getFeatureActiveLevel(&cluster, nil, &params, nil)).To(Equal(activeLevelActive))
				Expect((&ClusterManagedNetworkingFeature{}).getFeatureActiveLevel(&cluster, nil, &params, nil)).To(Equal(activeLevelActive))
			})
			It("ppc supporting minimal-iso", func() {
				cpuArchitecture := models.ClusterCPUArchitecturePpc64le
				cluster := common.Cluster{Cluster: models.Cluster{
					OpenshiftVersion: "4.12",
					CPUArchitecture:  cpuArchitecture,
				}}
				infraEnv := models.InfraEnv{Type: common.ImageTypePtr(models.ImageTypeMinimalIso)}
				Expect((&MinimalIso{}).getFeatureActiveLevel(&cluster, &infraEnv, nil, nil)).To(Equal(activeLevelActive))

				filters := SupportLevelFilters{OpenshiftVersion: "4.12", CPUArchitecture: &cpuArchitecture}
				supportLevel, _ := (&MinimalIso{}).getSupportLevel(filters)
				Expect(supportLevel).To(Equal(models.SupportLevelSupported))
			})

			for _, filters := range getPlatformFilters() {
				for _, feature := range []SupportLevelFeature{
					&VsphereIntegrationFeature{},
					&NutanixIntegrationFeature{},
					&BaremetalPlatformFeature{},
					&NonePlatformFeature{},
					&OciIntegrationFeature{},
					&ExternalPlatformFeature{},
				} {
					filters := filters
					feature := feature
					When("Empty support level - platforms", func() {
						It(fmt.Sprintf("Feature %s Platform %s ExternalPlatformName %s", feature.GetName(), *filters.PlatformType, swag.StringValue(filters.ExternalPlatformName)), func() {
							emptyFilters := SupportLevelFilters{OpenshiftVersion: "", CPUArchitecture: nil, PlatformType: nil, ExternalPlatformName: nil}
							supportLevel, _ := feature.getSupportLevel(emptyFilters)
							Expect(string(supportLevel)).To(Not(Equal("")))

							supportLevel, _ = feature.getSupportLevel(filters)
							Expect(string(supportLevel)).To(Equal(""))
						})
					})
				}
			}

			for _, filters := range getPlatformFilters() {
				filters := filters
				When("Empty support level - PlatformManagedNetworkingFeature", func() {
					It(string(*filters.PlatformType)+" "+swag.StringValue(filters.ExternalPlatformName), func() {
						feature := &PlatformManagedNetworkingFeature{}

						emptyFilters := SupportLevelFilters{OpenshiftVersion: "", CPUArchitecture: nil, PlatformType: nil, ExternalPlatformName: nil}
						supportLevel, _ := feature.getSupportLevel(emptyFilters)
						Expect(string(supportLevel)).To(Equal(""))

						supportLevel, _ = feature.getSupportLevel(filters)
						Expect(string(supportLevel)).To(Not(Equal("")))
					})
				})
			}
			It("s390x not supporting minimal-iso", func() {
				cpuArchitecture := models.ClusterCPUArchitectureS390x
				cluster := common.Cluster{Cluster: models.Cluster{
					OpenshiftVersion: "4.12",
					CPUArchitecture:  cpuArchitecture,
				}}
				infraEnv := models.InfraEnv{Type: common.ImageTypePtr(models.ImageTypeMinimalIso)}
				Expect((&MinimalIso{}).getFeatureActiveLevel(&cluster, &infraEnv, nil, nil)).To(Equal(activeLevelActive))

				filters := SupportLevelFilters{OpenshiftVersion: "4.12", CPUArchitecture: &cpuArchitecture}
				supportLevel, _ := (&MinimalIso{}).getSupportLevel(filters)
				Expect(supportLevel).To(Equal(models.SupportLevelUnavailable))
			})

			It("s390x not supporting minimal-iso without cluster", func() {
				cpuArchitecture := models.ClusterCPUArchitectureS390x

				infraEnv := models.InfraEnv{Type: common.ImageTypePtr(models.ImageTypeMinimalIso)}
				Expect((&MinimalIso{}).getFeatureActiveLevel(nil, &infraEnv, nil, nil)).To(Equal(activeLevelActive))

				filters := SupportLevelFilters{OpenshiftVersion: "", CPUArchitecture: &cpuArchitecture}
				supportLevel, _ := (&MinimalIso{}).getSupportLevel(filters)
				Expect(supportLevel).To(Equal(models.SupportLevelUnavailable))
			})

			It("Disable olm operator activated features in cluster", func() {
				operators := []*models.MonitoredOperator{
					{
						Name:             "cnv",
						Namespace:        "openshift-cnv",
						OperatorType:     models.OperatorTypeOlm,
						SubscriptionName: "hco-operatorhub",
						TimeoutSeconds:   60 * 60,
					},
				}

				cluster := common.Cluster{Cluster: models.Cluster{
					OpenshiftVersion:   "4.8",
					CPUArchitecture:    models.ClusterCPUArchitecturePpc64le,
					MonitoredOperators: operators,
				}}
				params := models.V2ClusterUpdateParams{
					OlmOperators: []*models.OperatorCreateParams{},
				}

				Expect((&CnvFeature{}).getFeatureActiveLevel(&cluster, nil, nil, nil)).To(Equal(activeLevelActive))
				Expect((&CnvFeature{}).getFeatureActiveLevel(&cluster, nil, &params, nil)).To(Equal(activeLevelNotActive))
			})
		})

		Context("getSupportLevel", func() {
			It("featuressupport.getSupportLevel equal to Feature.getSupportLevel", func() {
				featureA := ClusterManagedNetworkingFeature{}
				openshiftVersion := "4.13"
				cpuArchitecture := models.ClusterCPUArchitectureS390x
				filters := SupportLevelFilters{OpenshiftVersion: openshiftVersion, CPUArchitecture: &cpuArchitecture}
				supportLevel, _ := featureA.getSupportLevel(filters)
				equalSupportLevel := GetSupportLevel(featureA.getId(), filters)
				Expect(supportLevel).To(Equal(equalSupportLevel))
			})
		})

		Context("getIncompatibleFeatures", func() {
			It("Features without any restrictions", func() {
				features := []models.FeatureSupportLevelID{
					models.FeatureSupportLevelIDCUSTOMMANIFEST,
					models.FeatureSupportLevelIDSINGLENODEEXPANSION,
				}
				for _, featureId := range features {
					Expect(featuresList[featureId].getIncompatibleFeatures("")).To(BeNil())
				}
			})

			It("incompatibleFeatures - all features - no openshift version", func() {
				for featureId, feature := range featuresList {
					featureId := featureId
					feature := feature

					incompatibleFeatures := feature.getIncompatibleFeatures("")
					for _, incompatibleFeatureId := range incompatibleFeatures {
						incompatibleFeature := featuresList[incompatibleFeatureId]
						By(fmt.Sprintf("Feature  %s with incompatible feature %s", featureId, incompatibleFeatureId), func() {
							incompatibleFeatures2 := incompatibleFeature.getIncompatibleFeatures("")
							if incompatibleFeatures2 == nil {
								incompatibleFeatures2 = []models.FeatureSupportLevelID{}
							}
							Expect(incompatibleFeatures2).To(ContainElement(featureId))
						})
					}

				}
			})

			It("vSphere with dual-stack", func() {
				dualStackFeature := featuresList[models.FeatureSupportLevelIDDUALSTACK]
				vsphereFeature := featuresList[models.FeatureSupportLevelIDVSPHEREINTEGRATION]

				isDualStackIncompatibleWithVsphere := isFeatureCompatible("4.8", dualStackFeature, vsphereFeature)
				isVsphereIncompatibleWithDualStack := isFeatureCompatible("4.8", vsphereFeature, dualStackFeature)
				Expect((*isDualStackIncompatibleWithVsphere).getId()).To(Equal(vsphereFeature.getId()))
				Expect((*isVsphereIncompatibleWithDualStack).getId()).To(Equal(dualStackFeature.getId()))

				isDualStackIncompatibleWithVsphere = isFeatureCompatible("4.13", dualStackFeature, vsphereFeature)
				isVsphereIncompatibleWithDualStack = isFeatureCompatible("4.13", vsphereFeature, dualStackFeature)
				Expect(isDualStackIncompatibleWithVsphere).To(BeNil())
				Expect(isVsphereIncompatibleWithDualStack).To(BeNil())
			})
		})

		Context("Test validate active features", func() {
			It("Ignore validation on AddHostCluster", func() {
				logBuffer := bytes.Buffer{}
				testLogger := logrus.New()
				testLogger.SetOutput(&logBuffer)

				clusterID := strfmt.UUID("e679ea3f-3b85-40e0-8dc9-82fd6945d9b2")
				cluster := common.Cluster{Cluster: models.Cluster{
					ID:                &clusterID,
					OpenshiftVersion:  "4.19",
					Kind:              swag.String(models.ClusterKindAddHostsCluster),
					CPUArchitecture:   models.ClusterCPUArchitectureS390x,
					ControlPlaneCount: common.MinMasterHostsNeededForInstallationInHaMode,
					Platform:          &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNone)},
				}}

				infraEnv := models.InfraEnv{
					ClusterID:       *cluster.ID,
					CPUArchitecture: models.ClusterCPUArchitectureS390x,
				}

				Expect(ValidateIncompatibleFeatures(logrus.NewEntry(testLogger), models.ClusterCPUArchitectureS390x, &cluster, &infraEnv, nil)).To(BeNil())
				Expect(logBuffer.String()).To(ContainSubstring("skipping feature support validation:"))
			})

			DescribeTable(
				"Valid VipDhcpAllocation and OpenShift version",
				func(openshiftVersion string) {
					Expect(IsFeatureAvailable(models.FeatureSupportLevelIDVIPAUTOALLOC, openshiftVersion, swag.String("anyarch"))).To(BeFalse())
				},
				Entry("VipAutoAllocation disabled for 4.15", "4.15.3"),
				Entry("VipAutoAllocation disabled for 4.16", "4.16.2"),
			)

			DescribeTable(
				"Valid VipDhcpAllocation and OpenShift version",
				func(openshiftVersion string) {
					Expect(IsFeatureAvailable(models.FeatureSupportLevelIDVIPAUTOALLOC, openshiftVersion, swag.String("anyarch"))).To(BeTrue())
				},
				Entry("VipAutoAllocation enabled for 4.14", "4.14.3"),
				Entry("VipAutoAllocation enabled for 4.12", "4.12.24"),
			)

			DescribeTable(
				"Valid Network Type and OpenShift version",
				func(openshiftVersion, networkType string) {
					cluster := &common.Cluster{Cluster: models.Cluster{
						OpenshiftVersion: openshiftVersion,
						NetworkType:      &networkType,
						Platform:         &models.Platform{Type: common.PlatformTypePtr(models.PlatformTypeNutanix)},
					}}
					log := logrus.New()

					err := ValidateActiveFeatures(log, cluster, nil, nil)
					Expect(err).ShouldNot(HaveOccurred())
				},
				Entry("SDN Active with Openshift < 4.15", "4.14.3", models.ClusterNetworkTypeOpenShiftSDN),
				Entry("OVN Active with Openshift < 4.15", "4.14.3", models.ClusterNetworkTypeOVNKubernetes),
				Entry("OVN Active with Openshift = 4.15", "4.15.2", models.ClusterNetworkTypeOVNKubernetes),
				Entry("OVN Active with Openshift > 4.15", "4.18.9", models.ClusterNetworkTypeOVNKubernetes),
			)

			DescribeTable(
				"Invalid Network Type and OpenShift version",
				func(openshiftVersion, networkType, expectedErrorMessage string) {
					cluster := &common.Cluster{Cluster: models.Cluster{
						OpenshiftVersion: openshiftVersion,
						NetworkType:      &networkType,
					}}
					log := logrus.New()

					err := ValidateActiveFeatures(log, cluster, nil, nil)
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal(expectedErrorMessage))
				},
				Entry("SDN Active with Openshift 4.15", "4.15.3", models.ClusterNetworkTypeOpenShiftSDN, "Openshift version 4.15.3 is not supported for OpenShiftSDN NetworkType"),
				Entry("SDN Active with Openshift > 4.15", "4.18.3", models.ClusterNetworkTypeOpenShiftSDN, "Openshift version 4.18.3 is not supported for OpenShiftSDN NetworkType"),
			)

		})
	})

	DescribeTable(
		"User Managed load balancer support",
		func(version string, expected bool) {
			arch := "x86_64"
			actual := IsFeatureAvailable(
				models.FeatureSupportLevelIDUSERMANAGEDLOADBALANCER,
				version,
				&arch,
			)
			Expect(actual).To(Equal(expected))
		},
		Entry(
			"Not in 4.14",
			"4.14",
			false,
		),
		Entry(
			"Not in 4.15",
			"4.14",
			false,
		),
		Entry(
			"Yes in 4.16",
			"4.16",
			true,
		),
		Entry(
			"Yes in 4.17",
			"4.17",
			true,
		),
	)

	DescribeTable("User Managed load balancer compatability with other features", func(activeFeatures []SupportLevelFeature, shouldSucceed bool) {
		activeFeatures = append(activeFeatures, &UserManagedLoadBalancerFeature{})

		if shouldSucceed {
			Expect(
				isFeaturesCompatibleWithFeatures(
					common.MinimumVersionForUserManagedLoadBalancerFeature,
					activeFeatures),
			).ToNot(HaveOccurred())
		} else {
			Expect(
				isFeaturesCompatibleWithFeatures(
					common.MinimumVersionForUserManagedLoadBalancerFeature,
					activeFeatures),
			).To(HaveOccurred())
		}
	},
		Entry(
			"platform baremetal",
			[]SupportLevelFeature{&BaremetalPlatformFeature{}},
			true,
		),

		Entry(
			"external platform",
			[]SupportLevelFeature{&ExternalPlatformFeature{}},
			false,
		),

		Entry(
			"nutanix platform",
			[]SupportLevelFeature{&NutanixIntegrationFeature{}},
			false,
		),

		Entry(
			"vsphere platform",
			[]SupportLevelFeature{&VsphereIntegrationFeature{}},
			true,
		),

		Entry(
			"none platform",
			[]SupportLevelFeature{&NonePlatformFeature{}},
			false,
		),

		Entry(
			"user managed networking",
			[]SupportLevelFeature{&UserManagedNetworkingFeature{}},
			false,
		),

		Entry(
			"SNO",
			[]SupportLevelFeature{&UserManagedNetworkingFeature{}},
			false,
		),

		Entry(
			"dual stack",
			[]SupportLevelFeature{&DualStackFeature{}},
			false,
		),

		Entry(
			"dual stack vips",
			[]SupportLevelFeature{&DualStackVipsFeature{}},
			false,
		),

		Entry(
			"vip automatic allocation",
			[]SupportLevelFeature{&VipAutoAllocFeature{}},
			false,
		),
	)

})

func TestOperators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feature-Support-Level tests")
}
