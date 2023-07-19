// this file is copied from mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.*
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class packageDetails(name: String, displayName: String, environment: String) {
    val packageName = name
    val displayName = displayName
    val environment = environment

    fun buildConfiguration(providerName : String, path : String, manualVcsRoot: AbsoluteId, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, daysOfWeek: String, daysOfMonth: String) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                ConfigureGoEnv()
                DownloadTerraformBinary()
                // Adds steps:
                // - Determine Working Directory for this Package
                // - Pre-Sweeper
                // - Compile Test Binary
                // - Run via jen20/teamcity-go-test
                // - Post-Sweeper
                RunAcceptanceTests(path, packageName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(parallelism, "TestAcc", "12", "us-central1", "")
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path, packageName)
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour, daysOfWeek, daysOfMonth)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_SERVICE_%s_%s".format(provider.replace("-", "").toUpperCase(), environment.toUpperCase(), packageName.toUpperCase())
    }
}
