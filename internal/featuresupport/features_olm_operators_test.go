package featuresupport

import (
	"fmt"

	"github.com/go-openapi/swag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
)

var _ = Describe("V2ListFeatureSupportLevels API", func() {
	featureCNV := models.FeatureSupportLevelIDCNV
	lVMavailableVersions := []string{"4.11", "4.12", "4.13", "4.14", "4.15"}
	unspportedLVMVersions := []string{"4.10", "4.9", "4.8", "4.7", "4.6"}

	Context("Test LVM/Nutanix are not supported under 4.11", func() {
		features := []models.FeatureSupportLevelID{models.FeatureSupportLevelIDLVM, models.FeatureSupportLevelIDNUTANIXINTEGRATION}
		for _, f := range features {
			feature := f
			It(fmt.Sprintf("%s test", feature), func() {
				for _, version := range unspportedLVMVersions {
					Expect(IsFeatureAvailable(feature, version, nil)).To(BeFalse())
				}
				for _, version := range lVMavailableVersions {
					Expect(IsFeatureAvailable(feature, version, nil)).To(BeTrue())
				}
				// feature test
				Expect(IsFeatureAvailable(feature, "4.30", nil)).To(BeTrue())

			})
		}
	})

	Context("Test LVM feature", func() {
		lvmFeatureList := featuresList[models.FeatureSupportLevelIDLVM]
		feature := models.FeatureSupportLevelIDLVM
		It("Validate LVM on CPU arch", func() {
			supportedCpuArch := []string{
				models.ClusterCPUArchitectureArm64,
				models.ClusterCPUArchitectureMulti,
				models.ClusterCPUArchitectureX8664,
			}
			notSupportedCpuArch := []string{
				models.ClusterCPUArchitectureS390x,
				models.ClusterCPUArchitecturePpc64le,
			}
			for _, arch := range supportedCpuArch {
				Expect(IsFeatureAvailable(feature, "4.11", swag.String(arch))).To(BeTrue())
			}
			for _, arch := range notSupportedCpuArch {
				Expect(IsFeatureAvailable(feature, "4.11", swag.String(arch))).To(BeFalse())
			}
		})
		It("Validate Feature Support for LVM", func() {

			tests := []struct {
				id       int // used to know which test case failed
				version  string
				platform *models.PlatformType
				expected models.SupportLevel
			}{
				{
					id:       1,
					version:  "4.11",
					platform: models.PlatformTypeNone.Pointer(),
					expected: models.SupportLevelDevPreview,
				},
				{
					id:       2,
					version:  "4.9",
					platform: models.PlatformTypeBaremetal.Pointer(),
					expected: models.SupportLevelUnavailable,
				},
				{
					id:       3,
					version:  "4.11",
					platform: models.PlatformTypeVsphere.Pointer(),
					expected: models.SupportLevelUnavailable,
				},
				{
					id:       4,
					version:  "4.12",
					platform: models.PlatformTypeBaremetal.Pointer(),
					expected: models.SupportLevelSupported,
				},
				{
					id:       5,
					version:  "4.14",
					platform: models.PlatformTypeNone.Pointer(),
					expected: models.SupportLevelSupported,
				},
				{
					id:       6,
					version:  "4.15",
					platform: models.PlatformTypeNone.Pointer(),
					expected: models.SupportLevelSupported,
				},
			}

			for _, test := range tests {

				featureSupportParams := SupportLevelFilters{OpenshiftVersion: test.version, CPUArchitecture: nil, PlatformType: test.platform}
				resultSupportLevel := GetSupportLevel(feature, featureSupportParams)
				Expect(fmt.Sprintf("id: %d, result: %s", test.id, resultSupportLevel)).To(Equal(fmt.Sprintf("id: %d, result: %s", test.id, test.expected)))
			}
		})
		It("Validate Incompatible Features", func() {
			incompatibleFeatures := make(map[string][]models.FeatureSupportLevelID)

			incompatibleFeatures["4.11"] = []models.FeatureSupportLevelID{
				models.FeatureSupportLevelIDNUTANIXINTEGRATION,
				models.FeatureSupportLevelIDVSPHEREINTEGRATION,
				models.FeatureSupportLevelIDODF,
				models.FeatureSupportLevelIDOPENSHIFTAI,
				models.FeatureSupportLevelIDVIPAUTOALLOC,
				models.FeatureSupportLevelIDCLUSTERMANAGEDNETWORKING,
			}

			incompatibleFeatures["4.12"] = []models.FeatureSupportLevelID{
				models.FeatureSupportLevelIDNUTANIXINTEGRATION,
				models.FeatureSupportLevelIDVSPHEREINTEGRATION,
				models.FeatureSupportLevelIDODF,
				models.FeatureSupportLevelIDOPENSHIFTAI,
				models.FeatureSupportLevelIDVIPAUTOALLOC,
				models.FeatureSupportLevelIDCLUSTERMANAGEDNETWORKING,
			}

			incompatibleFeatures["4.15"] = []models.FeatureSupportLevelID{
				models.FeatureSupportLevelIDNUTANIXINTEGRATION,
				models.FeatureSupportLevelIDVSPHEREINTEGRATION,
				models.FeatureSupportLevelIDODF,
				models.FeatureSupportLevelIDOPENSHIFTAI,
			}
			incompatibleFeatures["4.16.0-rc0"] = []models.FeatureSupportLevelID{
				models.FeatureSupportLevelIDNUTANIXINTEGRATION,
				models.FeatureSupportLevelIDVSPHEREINTEGRATION,
				models.FeatureSupportLevelIDODF,
				models.FeatureSupportLevelIDOPENSHIFTAI,
			}

			testIncompatibleFeatures := []struct {
				id          int
				version     string
				featureList []models.FeatureSupportLevelID
			}{
				{
					id:          1,
					version:     "4.11",
					featureList: []models.FeatureSupportLevelID{models.FeatureSupportLevelIDLVM},
				},
				{
					id:          2,
					version:     "4.12",
					featureList: []models.FeatureSupportLevelID{models.FeatureSupportLevelIDLVM},
				},
				{
					id:          3,
					version:     "4.15",
					featureList: []models.FeatureSupportLevelID{models.FeatureSupportLevelIDLVM},
				},
				{
					id:          4,
					version:     "4.16.0-rc0", // check pre release version
					featureList: []models.FeatureSupportLevelID{models.FeatureSupportLevelIDLVM},
				},
			}

			for _, test := range testIncompatibleFeatures {
				for _, featureId := range test.featureList {
					result := featuresList[featureId].getIncompatibleFeatures(test.version)
					Expect(fmt.Sprintf("id: %d, result: %s", test.id, result)).To(Equal(fmt.Sprintf("id: %d, result: %s", test.id, incompatibleFeatures[test.version])))
				}
			}
		})

		It("Ensure LVM  multinode is supportted on version 4.15", func() {
			features := []models.FeatureSupportLevelID{
				models.FeatureSupportLevelIDCLUSTERMANAGEDNETWORKING,
				models.FeatureSupportLevelIDVIPAUTOALLOC,
				models.FeatureSupportLevelIDSNO,
			}
			for _, feature := range features {
				Expect(isFeatureCompatible("4.15", featuresList[feature], lvmFeatureList)).To(BeNil())
			}
		})
	})

	Context("Test feature support levels for Nutanix platform", func() {
		It("CNV should be unavailable", func() {
			featureSupportLevels := GetFeatureSupportList(
				"4.14",
				swag.String(common.X86CPUArchitecture),
				common.PlatformTypePtr(models.PlatformTypeNutanix),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDCNV {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelUnavailable))
				Expect(feature.Reason).To(Equal(models.IncompatibilityReasonPlatform))
			}
		})

		It("MCE should be unavailable", func() {
			featureSupportLevels := GetFeatureSupportList(
				"4.14",
				swag.String(common.X86CPUArchitecture),
				common.PlatformTypePtr(models.PlatformTypeNutanix),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDMCE {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelUnavailable))
				Expect(feature.Reason).To(Equal(models.IncompatibilityReasonPlatform))
			}
		})
	})

	Context("Test feature support levels for Vsphere platform", func() {
		It("CNV should be unavailable", func() {
			featureSupportLevels := GetFeatureSupportList(
				"4.14",
				swag.String(common.X86CPUArchitecture),
				common.PlatformTypePtr(models.PlatformTypeVsphere),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDCNV {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(models.SupportLevelUnavailable))
				Expect(feature.Reason).To(Equal(models.IncompatibilityReasonPlatform))
			}
		})
	})

	DescribeTable("Validate CNV on Architecture", func(ocpVersion []string, cpuArch string, expectedResult bool) {
		for _, v := range ocpVersion {
			version := v
			result := IsFeatureAvailable(featureCNV, version, swag.String(cpuArch))
			Expect(result).Should(Equal(expectedResult),
				fmt.Sprintf("Feature: %s, OCP version: %s, CpuArch: %s, should be %v", featureCNV, v, cpuArch, expectedResult))
		}
	},

		Entry("on X86	is supported", []string{"4.8", "4.11", "4.14", "4.21"}, models.ClusterCPUArchitectureX8664, true),
		Entry("on arm64 is supported", []string{"4.14", "4.15", "4.21"}, models.ClusterCPUArchitectureArm64, true),
		Entry("on arm64 is    NOT supported", []string{"4.11", "4.12", "4.13"}, models.ClusterCPUArchitectureArm64, false),
		Entry("on S390x is    NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitectureS390x, false),
		Entry("on ppc64le is	NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitecturePpc64le, false),
	)

	Context("Test MTV feature", func() {
		DescribeTable("Validate MTV on Architecture", func(ocpVersion []string, cpuArch string, expectedResult bool) {
			for _, v := range ocpVersion {
				version := v
				result := IsFeatureAvailable(models.FeatureSupportLevelIDMTV, version, swag.String(cpuArch))
				Expect(result).Should(Equal(expectedResult),
					fmt.Sprintf("Feature: %s, OCP version: %s, CpuArch: %s, should be %v", models.FeatureSupportLevelIDMTV, v, cpuArch, expectedResult))
			}
		},
			Entry("on X86	is supported above 4.14", []string{"4.14", "4.21"}, models.ClusterCPUArchitectureX8664, true),
			Entry("on X86	is NOT supported", []string{"4.8", "4.11", "4.13"}, models.ClusterCPUArchitectureX8664, false),
			Entry("on arm64 is    NOT supported", []string{"4.8", "4.11", "4.14", "4.21"}, models.ClusterCPUArchitectureArm64, false),
			Entry("on S390x is    NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitectureS390x, false),
			Entry("on ppc64le is	NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitecturePpc64le, false),
		)

		DescribeTable("Validate MTV on platform", func(ocpVersion string, cpuArch string, platformType models.PlatformType, expectedResult models.SupportLevel) {
			featureSupportLevels := GetFeatureSupportList(
				ocpVersion,
				swag.String(cpuArch),
				common.PlatformTypePtr(platformType),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDMTV {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(expectedResult))
			}
		},
			Entry("on Vsphere", "4.14", common.X86CPUArchitecture, models.PlatformTypeVsphere, models.SupportLevelUnavailable),
			Entry("on Nutanix", "4.14", common.X86CPUArchitecture, models.PlatformTypeNutanix, models.SupportLevelUnavailable),
			Entry("on none", "4.14", common.X86CPUArchitecture, models.PlatformTypeNone, models.SupportLevelSupported),
			Entry("on baremetal", "4.14", common.X86CPUArchitecture, models.PlatformTypeBaremetal, models.SupportLevelSupported),
			Entry("on baremetal", "4.13", common.X86CPUArchitecture, models.PlatformTypeBaremetal, models.SupportLevelUnavailable),
		)
	})

	Context("Test OSC feature", func() {
		DescribeTable("Validate OSC on Architecture", func(ocpVersion []string, cpuArch string, expectedResult bool) {
			for _, v := range ocpVersion {
				version := v
				result := IsFeatureAvailable(models.FeatureSupportLevelIDOSC, version, swag.String(cpuArch))
				Expect(result).Should(Equal(expectedResult),
					fmt.Sprintf("Feature: %s, OCP version: %s, CpuArch: %s, should be %v", models.FeatureSupportLevelIDOSC, v, cpuArch, expectedResult))
			}
		},
			Entry("on X86	is supported above 4.10", []string{"4.15", "4.16", "4.17", "4.21"}, models.ClusterCPUArchitectureX8664, true),
			Entry("on arm64 is    NOT supported", []string{"4.8", "4.11", "4.14", "4.21"}, models.ClusterCPUArchitectureArm64, false),
			Entry("on S390x is    NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitectureS390x, false),
			Entry("on ppc64le is	NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitecturePpc64le, false),
		)

		DescribeTable("Validate OSC on platform", func(ocpVersion string, cpuArch string, platformType models.PlatformType, expectedResult models.SupportLevel) {
			featureSupportLevels := GetFeatureSupportList(
				ocpVersion,
				swag.String(cpuArch),
				common.PlatformTypePtr(platformType),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDOSC {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(expectedResult))
			}
		},
			Entry("on Vsphere", "4.10", common.X86CPUArchitecture, models.PlatformTypeVsphere, models.SupportLevelUnavailable),
			Entry("on Nutanix", "4.10", common.X86CPUArchitecture, models.PlatformTypeNutanix, models.SupportLevelUnavailable),
			Entry("on none", "4.10", common.X86CPUArchitecture, models.PlatformTypeNone, models.SupportLevelTechPreview),
			Entry("on baremetal", "4.10", common.X86CPUArchitecture, models.PlatformTypeBaremetal, models.SupportLevelTechPreview),
		)
	})

	Context("Test NMState feature", func() {
		DescribeTable("Validate NMState on Architecture", func(ocpVersion []string, cpuArch string, expectedResult bool) {
			for _, v := range ocpVersion {
				version := v
				result := IsFeatureAvailable(models.FeatureSupportLevelIDNMSTATE, version, swag.String(cpuArch))
				Expect(result).Should(Equal(expectedResult),
					fmt.Sprintf("Feature: %s, OCP version: %s, CpuArch: %s, should be %v", models.FeatureSupportLevelIDNMSTATE, v, cpuArch, expectedResult))
			}
		},
			Entry("on X86	is supported above 4.12", []string{"4.12", "4.16", "4.17", "4.21"}, models.ClusterCPUArchitectureX8664, true),
			Entry("on X86 is 	NOT supported below 4.12", []string{"3.18", "4.1", "4.10", "4.11"}, models.ClusterCPUArchitectureX8664, false),
			Entry("on arm64 is    NOT supported", []string{"4.8", "4.11", "4.14", "4.21"}, models.ClusterCPUArchitectureArm64, false),
			Entry("on S390x is    NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitectureS390x, false),
			Entry("on ppc64le is	NOT supported", []string{"4.11", "4.13", "4.14", "4.21"}, models.ClusterCPUArchitecturePpc64le, false),
		)

		DescribeTable("Validate NMState on platform", func(ocpVersion string, cpuArch string, platformType models.PlatformType, expectedResult models.SupportLevel) {
			featureSupportLevels := GetFeatureSupportList(
				ocpVersion,
				swag.String(cpuArch),
				common.PlatformTypePtr(platformType),
				nil,
			)

			for _, feature := range featureSupportLevels {
				if feature.FeatureSupportLevelID != models.FeatureSupportLevelIDNMSTATE {
					continue
				}

				Expect(feature.SupportLevel).To(Equal(expectedResult))
			}
		},
			Entry("on Vsphere", "4.12", common.X86CPUArchitecture, models.PlatformTypeVsphere, models.SupportLevelSupported),
			Entry("on none", "4.12", common.X86CPUArchitecture, models.PlatformTypeNone, models.SupportLevelSupported),
			Entry("on baremetal", "4.12", common.X86CPUArchitecture, models.PlatformTypeBaremetal, models.SupportLevelSupported),

			Entry("on Nutanix", "4.12", common.X86CPUArchitecture, models.PlatformTypeNutanix, models.SupportLevelUnavailable),
			Entry("on Nutanix", "4.12", common.X86CPUArchitecture, models.PlatformTypeExternal, models.SupportLevelUnavailable),
		)
	})
})
