/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

import DefaultTerraformCoreVersion
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.ParameterDisplay
import jetbrains.buildServer.configs.kotlin.ParametrizedWithType

// this file is copied from mmv1, any changes made here will be overwritten

// NOTE: this file includes Extensions of the Kotlin DSL class ParametrizedWithType
// This allows us to reuse code in the config easily, while ensuring the same parameters are set across builds.
// See the class's documentation: https://teamcity.jetbrains.com/app/dsl-documentation/root/parametrized-with-type/index.html

// AllContextParameters is used to pass ALL the values set via Context Parameters into the Kotlin code's entrypoint.
class AllContextParameters(

    // Values that differ across environments

    // GOOGLE_CREDENTIALS
    val credentialsGa: String,
    val credentialsBeta: String,
    val credentialsVcr: String,

    // GOOGLE_SERVICE_ACCOUNT
    val serviceAccountGa: String,
    val serviceAccountBeta: String,
    val serviceAccountVcr: String,

    // GOOGLE_PROJECT & GOOGLE_PROJECT_NUMBER
    val projectGa: String,
    val projectBeta: String,
    val projectVcr: String,
    val projectNumberGa: String,
    val projectNumberBeta: String,
    val projectNumberVcr: String,

    // GOOGLE_IDENTITY_USER
    val identityUserGa: String,
    val identityUserBeta: String,
    val identityUserVcr: String,

    // GOOGLE_MASTER_BILLING_ACCOUNT
    val masterBillingAccountGa: String,
    val masterBillingAccountBeta: String,
    val masterBillingAccountVcr: String,

    // GOOGLE_ORG_2
    val org2Ga: String,
    val org2Beta: String,
    val org2Vcr: String,

    // GOOGLE_CHRONICLE_INSTANCE_ID
    val chronicleInstanceIdGa: String,
    val chronicleInstanceIdBeta: String,
    val chronicleInstanceIdVcr: String,

    // GOOGLE_VMWAREENGINE_PROJECT
    val vmwareengineProjectGa: String,
    val vmwareengineProjectBeta: String,
    val vmwareengineProjectVcr: String,

    // Values that are the same across GA, Beta, and VCR testing environments
    val billingAccount: String,   // GOOGLE_BILLING_ACCOUNT
    val billingAccount2: String,  // GOOGLE_BILLING_ACCOUNT_2
    val custId: String,           // GOOGLE_CUST_ID
    val org: String,              // GOOGLE_ORG
    val orgDomain: String,        // GOOGLE_ORG_DOMAIN
    val region: String,           // GOOGLE_REGION
    val zone: String,             // GOOGLE_ZONE

    // VCR specific
    val infraProject: String,     // GOOGLE_INFRA_PROJECT
    val vcrBucketName: String,    // VCR_BUCKET_NAME

    // GCS specific (for nightly + upstream MM logs)
    val credentialsGCS: String,   // GOOGLE_CREDENTIALS_GCS
    )

// AccTestConfiguration is used to easily pass values set via Context Parameters into build configurations.
class AccTestConfiguration(
    val billingAccount: String,
    val billingAccount2: String,
    val credentials: String,
    val custId: String,
    val identityUser: String,
    val masterBillingAccount: String,
    val org: String,
    val org2: String,
    val chronicleInstanceId: String,
    val orgDomain: String,
    val project: String,
    val projectNumber: String,
    val region: String,
    val serviceAccount: String,
    val vmwareengineProject: String,
    val zone: String,

    // VCR specific
    val infraProject: String,
    val vcrBucketName: String,

    // GCS specific (for nightly + upstream MM logs)
    val credentialsGCS: String,
    )

fun getGaAcceptanceTestConfig(allConfig: AllContextParameters): AccTestConfiguration {
    return AccTestConfiguration(
        allConfig.billingAccount,
        allConfig.billingAccount2,
        allConfig.credentialsGa,
        allConfig.custId,
        allConfig.identityUserGa,
        allConfig.masterBillingAccountGa,
        allConfig.org,
        allConfig.org2Ga,
        allConfig.chronicleInstanceIdGa,
        allConfig.orgDomain,
        allConfig.projectGa,
        allConfig.projectNumberGa,
        allConfig.region,
        allConfig.serviceAccountGa,
        allConfig.vmwareengineProjectGa,
        allConfig.zone,
        allConfig.infraProject,
        allConfig.vcrBucketName,
        allConfig.credentialsGCS
    )
}

fun getBetaAcceptanceTestConfig(allConfig: AllContextParameters): AccTestConfiguration {
    return AccTestConfiguration(
        allConfig.billingAccount,
        allConfig.billingAccount2,
        allConfig.credentialsBeta,
        allConfig.custId,
        allConfig.identityUserBeta,
        allConfig.masterBillingAccountBeta,
        allConfig.org,
        allConfig.org2Beta,
        allConfig.chronicleInstanceIdBeta,
        allConfig.orgDomain,
        allConfig.projectBeta,
        allConfig.projectNumberBeta,
        allConfig.region,
        allConfig.serviceAccountBeta,
        allConfig.vmwareengineProjectBeta,
        allConfig.zone,
        allConfig.infraProject,
        allConfig.vcrBucketName,
        allConfig.credentialsGCS
    )
}

fun getVcrAcceptanceTestConfig(allConfig: AllContextParameters): AccTestConfiguration {
    return AccTestConfiguration(
        allConfig.billingAccount,
        allConfig.billingAccount2,
        allConfig.credentialsVcr,
        allConfig.custId,
        allConfig.identityUserVcr,
        allConfig.masterBillingAccountVcr,
        allConfig.org,
        allConfig.org2Vcr,
        allConfig.chronicleInstanceIdVcr,
        allConfig.orgDomain,
        allConfig.projectVcr,
        allConfig.projectNumberVcr,
        allConfig.region,
        allConfig.serviceAccountVcr,
        allConfig.vmwareengineProjectVcr,
        allConfig.zone,
        allConfig.infraProject,
        allConfig.vcrBucketName,
        allConfig.credentialsGCS
    )
}

// ParametrizedWithType.configureGoogleSpecificTestParameters allows build configs to be created
// with the environment variables needed to configure the provider and/or configure test code.
fun ParametrizedWithType.configureGoogleSpecificTestParameters(config: AccTestConfiguration) {
    hiddenVariable("env.GOOGLE_BILLING_ACCOUNT", config.billingAccount, "The billing account associated with the first google organization")
    hiddenVariable("env.GOOGLE_BILLING_ACCOUNT_2", config.billingAccount2, "The billing account associated with the second google organization")
    hiddenVariable("env.GOOGLE_CUST_ID", config.custId, "The ID of the Google Identity Customer")
    hiddenVariable("env.GOOGLE_ORG", config.org, "The Google Organization Id")
    hiddenVariable("env.GOOGLE_ORG_2", config.org2, "The second Google Organization Id")
    hiddenVariable("env.GOOGLE_MASTER_BILLING_ACCOUNT", config.masterBillingAccount, "The master billing account")
    hiddenVariable("env.GOOGLE_PROJECT", config.project, "The google project for this build")
    hiddenVariable("env.GOOGLE_ORG_DOMAIN", config.orgDomain, "The org domain")
    hiddenVariable("env.GOOGLE_PROJECT_NUMBER", config.projectNumber, "The project number associated with the project")
    hiddenVariable("env.GOOGLE_REGION", config.region, "The google region to use")
    hiddenVariable("env.GOOGLE_SERVICE_ACCOUNT", config.serviceAccount, "The service account")
    hiddenVariable("env.GOOGLE_ZONE", config.zone, "The google zone to use")
    hiddenVariable("env.GOOGLE_IDENTITY_USER", config.identityUser, "The user for the identity platform")
    hiddenVariable("env.GOOGLE_CHRONICLE_INSTANCE_ID", config.chronicleInstanceId, "The id of the Chronicle instance")
    hiddenVariable("env.GOOGLE_VMWAREENGINE_PROJECT", config.vmwareengineProject, "The project used for vmwareengine tests")
    hiddenPasswordVariable("env.GOOGLE_CREDENTIALS", config.credentials, "The Google credentials for this test runner")
}

// ParametrizedWithType.acceptanceTestBuildParams sets build params that affect how commands to run
//  acceptance tests are templated
fun ParametrizedWithType.acceptanceTestBuildParams(parallelism: Int, prefix: String, timeout: String, releaseDiffTest: Boolean) {
    hiddenVariable("env.TF_ACC", "1", "Set to a value to run the Acceptance Tests")
    text("PARALLELISM", "%d".format(parallelism))
    text("TEST_PREFIX", prefix)
    text("TIMEOUT", timeout)
    if (releaseDiffTest) {
        text("env.RELEASE_DIFF", "true")
    } else {
        // Use an empty string for backwards-compatibility with pre-7.X release diff behavior.
        text("env.RELEASE_DIFF", "")
    }
}

// ParametrizedWithType.sweeperParameters sets build parameters that affect how sweepers are run
fun ParametrizedWithType.sweeperParameters(sweeperRegions: String, sweepRun: String) {
    text("SWEEPER_REGIONS", sweeperRegions)
    text("SWEEP_RUN", sweepRun)
}

// ParametrizedWithType.terraformSkipProjectSweeper sets an environment variable to skip the sweeper for project resources
fun ParametrizedWithType.terraformSkipProjectSweeper() {
    text("env.SKIP_PROJECT_SWEEPER", "1")
}

// BuildType.disableProjectSweep disabled sweeping project resources after a build configuration has been initialised
fun BuildType.disableProjectSweep(){
    params {
        terraformSkipProjectSweeper()
    }
}

// ParametrizedWithType.terraformEnableProjectSweeper unsets an environment variable used to skip the sweeper for project resources
fun ParametrizedWithType.terraformEnableProjectSweeper() {
    text("env.SKIP_PROJECT_SWEEPER", "")
}

// BuildType.enableProjectSweep enables sweeping project resources after a build configuration has been initialised
fun BuildType.enableProjectSweep(){
    params {
        terraformEnableProjectSweeper()
    }
}

// ParametrizedWithType.vcrEnvironmentVariables sets environment variables that influence how VCR tests run
// These values can be changed in custom builds, e.g. setting VCR_MODE=REPLAYING
fun ParametrizedWithType.vcrEnvironmentVariables(config: AccTestConfiguration, providerName: String) {
    text("env.VCR_MODE", "RECORDING")
    text("env.VCR_PATH", "%system.teamcity.build.checkoutDir%/fixtures")
    text("env.TEST", "./${providerName}/services/...")
    text("env.TESTARGS", "-run=%TEST_PREFIX%")
    hiddenVariable("env.GOOGLE_INFRA_PROJECT", config.infraProject, "The project that's linked to the GCS bucket storing VCR cassettes")
    hiddenVariable("env.VCR_BUCKET_NAME", config.vcrBucketName, "The name of the GCS bucket storing VCR cassettes")
}

// ParametrizedWithType.terraformLoggingParameters sets environment variables and build parameters that
// affect which logs are shown and allows them to be saved
fun ParametrizedWithType.terraformLoggingParameters(config: AccTestConfiguration, providerName: String) {
    // Set logging levels to match old projects
    text("env.TF_LOG", "DEBUG")
    text("env.TF_LOG_CORE", "WARN")
    text("env.TF_LOG_SDK_FRAMEWORK", "INFO")

    // Set where logs are sent
    text("PROVIDER_NAME", providerName)
    text("env.TF_LOG_PATH_MASK", "%system.teamcity.build.checkoutDir%/debug-%PROVIDER_NAME%-%env.BUILD_NUMBER%-%teamcity.build.id%-%s.txt") // .txt extension used to make artifacts open in browser, instead of download

    hiddenPasswordVariable("env.GOOGLE_CREDENTIALS_GCS", config.credentialsGCS, "The Google credentials for copying debug logs to the GCS bucket")
}

fun ParametrizedWithType.readOnlySettings() {
    hiddenVariable("teamcity.ui.settings.readOnly", "true", "Requires build configurations be edited via Kotlin")
}

// ParametrizedWithType.terraformCoreBinaryTesting sets environment variables that control what Terraform version is downloaded
// and ensures the testing framework uses that downloaded version. The default Terraform core version is used if no argument is supplied.
fun ParametrizedWithType.terraformCoreBinaryTesting(tfVersion: String = DefaultTerraformCoreVersion) {
    text("env.TERRAFORM_CORE_VERSION", tfVersion, "The version of Terraform Core which should be used for testing")
    hiddenVariable("env.TF_ACC_TERRAFORM_PATH", "%system.teamcity.build.checkoutDir%/tools/terraform", "The path where the Terraform Binary is located. Used by the testing framework.")
}

// BuildType.overrideTerraformCoreVersion is used to override the value of TERRAFORM_CORE_VERSION in special cases where we're testing new features
// that rely on a specific version of Terraform we might not want to be used for all our tests in TeamCity.
fun BuildType.overrideTerraformCoreVersion(tfVersion: String){
    params {
        terraformCoreBinaryTesting(tfVersion)
    }
}

fun ParametrizedWithType.terraformShouldPanicForSchemaErrors() {
    hiddenVariable("env.TF_SCHEMA_PANIC_ON_ERROR", "1", "Panic if unknown/unmatched fields are set into the state")
}

fun ParametrizedWithType.workingDirectory(path: String) {
    text("PACKAGE_PATH", path, "", "The path at which to run - automatically updated", ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}