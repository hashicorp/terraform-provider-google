/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.AbsoluteId

const val providerName = "google"

// Google returns an instance of Project,
// which has multiple build configurations defined within it.
// See https://teamcity.jetbrains.com/app/dsl-documentation/root/project/index.html
fun Google(environment: String, manualVcsRoot: AbsoluteId, branchRef: String, configuration: ClientConfiguration) : Project {
    return Project{

        // Create build configs for each package defined in packages.kt
        val packageConfigs = buildConfigurationsForPackages(packages, providerName, environment, manualVcsRoot, branchRef, configuration)
        packageConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }

        // Create build configs for each service package defined in services.kt
        val servicePackageConfigs = buildConfigurationsForPackages(services, providerName, environment, manualVcsRoot, branchRef, configuration)
        servicePackageConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }

        // Create build configs for sweepers, including dependencies on all builds made above
        val allDependencyIds = ArrayList<String>()
        (packageConfigs + servicePackageConfigs).forEach { config ->
            allDependencyIds.add(config.id.toString())
        }
        val sweeperPackageConfigs = buildConfigurationsForSweepers(sweepers, providerName, environment, manualVcsRoot, branchRef, configuration, allDependencyIds)
        sweeperPackageConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
    }
}

fun buildConfigurationsForPackages(packages: Map<String, Map<String, String>>, providerName : String, environment: String, manualVcsRoot: AbsoluteId, branchRef: String, environmentVariables: ClientConfiguration): List<BuildType> {
    val triggerConfig = NightlyTriggerConfiguration(environment, branchRef)
    var list = ArrayList<BuildType>()

    // Create build configurations for all packages, except sweeper
    packages.forEach { (packageName, info) ->

        val path: String = info.getValue("path").toString()
        val name: String = info.getValue("name").toString()
        val displayName: String = info.getValue("displayName").toString()

        val pkg = packageDetails(packageName, displayName, providerName, environment)
        val buildConfig = pkg.buildConfiguration(path, manualVcsRoot, defaultParallelism, environmentVariables)

        list.add(buildConfig)
    }

    return list
}

fun buildConfigurationsForSweepers(packages: Map<String, Map<String, String>>, providerName : String, environment: String, manualVcsRoot: AbsoluteId, branchRef: String, environmentVariables: ClientConfiguration, dependencies: ArrayList<String> ): List<BuildType> {

    val triggerConfig = NightlyTriggerConfiguration(environment, branchRef)
    var list = ArrayList<BuildType>()

    val sweeperPackage : Map<String, String> = packages.getValue("sweeper")
    val sweeperPath : String = sweeperPackage.getValue("path")!!.toString()
    val s = sweeperBuildConfigs()

    // Pre-Sweeper
    val preSweeperConfig = s.preSweeperBuildConfig(sweeperPath, manualVcsRoot, defaultParallelism, environmentVariables)
    list.add(preSweeperConfig)

    // Post-Sweeper + dependencies + trigger
    val postSweeperConfig = s.postSweeperBuildConfig(sweeperPath, manualVcsRoot, defaultParallelism, triggerConfig, environmentVariables, dependencies)
    list.add(postSweeperConfig)

    return list
}

class NightlyTriggerConfiguration(environment: String, branchRef: String, nightlyTestsEnabled: Boolean = true, startHour: Int = defaultStartHour, daysOfWeek: String = defaultDaysOfWeek, daysOfMonth: String = defaultDaysOfMonth) {

    // Default values are used below unless
    // - alternatives passed in as arguments
    // - logic in `init` changes them based on environment
    var branchRef = branchRef
    var nightlyTestsEnabled = nightlyTestsEnabled
    var startHour = startHour
    var daysOfWeek = daysOfWeek
    var daysOfMonth = daysOfMonth

    init {
        // If the environment parameter is set to the value of MAJOR_RELEASE_TESTING, 
        // change the days of week to the day for v5.0.0 feature branch testing
        if (environment == MAJOR_RELEASE_TESTING) {
            this.daysOfWeek = "4" // Thursday for GA
        }
    }

}
