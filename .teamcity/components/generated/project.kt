/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.*

const val providerName = "google"

// Google returns an instance of Project,
// which has multiple build configurations defined within it.
// See https://teamcity.jetbrains.com/app/dsl-documentation/root/project/index.html
fun Google(environment: String, manualVcsRoot: AbsoluteId, branchRef: String, configuration: ClientConfiguration) : Project {

    // Create build configs for each package defined in packages.kt and services.kt files
    val allPackages = packages + services
    val packageConfigs = buildConfigurationsForPackages(allPackages, providerName, environment, manualVcsRoot, configuration)

    // Create build configs for sweepers
    val preSweeperConfig = buildConfigurationForSweeper("Pre-Sweeper", sweepers, providerName, manualVcsRoot, configuration)
    val postSweeperConfig = buildConfigurationForSweeper("Post-Sweeper", sweepers, providerName, manualVcsRoot, configuration)

    // Add trigger to last step of build chain (post-sweeper)
    val triggerConfig = NightlyTriggerConfiguration(environment, branchRef)
    postSweeperConfig.addTrigger(triggerConfig)
    
    return Project{

        // Register build configs in the project
        buildType(preSweeperConfig)
        packageConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
        buildType(postSweeperConfig)

        // Set up dependencies between builds using `sequential` block
        // Acc test builds run in parallel
        sequential {
            buildType(preSweeperConfig)

            parallel{
                packageConfigs.forEach { buildConfiguration ->
                    buildType(buildConfiguration)
                }
            }

            buildType(postSweeperConfig)
        }
    }
}

fun buildConfigurationsForPackages(packages: Map<String, Map<String, String>>, providerName: String, environment: String, manualVcsRoot: AbsoluteId, environmentVariables: ClientConfiguration): List<BuildType> {
    var list = ArrayList<BuildType>()

    // Create build configurations for all packages, except sweeper
    packages.forEach { (packageName, info) ->
        val path: String = info.getValue("path").toString()
        val name: String = info.getValue("name").toString()
        val displayName: String = info.getValue("displayName").toString()

        val pkg = packageDetails(packageName, displayName, providerName)
        val buildConfig = pkg.buildConfiguration(path, manualVcsRoot, defaultParallelism, environmentVariables)
        list.add(buildConfig)
    }

    return list
}

fun buildConfigurationForSweeper(sweeperName: String, packages: Map<String, Map<String, String>>, providerName: String, manualVcsRoot: AbsoluteId, environmentVariables: ClientConfiguration): BuildType {
    val sweeperPackage : Map<String, String> = packages.getValue("sweeper")
    val sweeperPath : String = sweeperPackage.getValue("path")!!.toString()
    val s = sweeperDetails()

    return s.sweeperBuildConfig(sweeperName, sweeperPath, providerName, manualVcsRoot, defaultParallelism, environmentVariables)
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
