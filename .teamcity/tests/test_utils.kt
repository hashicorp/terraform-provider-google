/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.AllContextParameters
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert.fail

const val gaProjectName = "Google"
const val betaProjectName = "Google Beta"
const val nightlyTestsProjectName = "Nightly Tests"
const val mmUpstreamProjectName = "Upstream MM Testing"
const val projectSweeperProjectName = "Project Sweeper"

fun testContextParameters(): AllContextParameters {
    return AllContextParameters(
        "credsGa",
        "credsBeta",
        "credsVcr",
        "serviceAccountGa",
        "serviceAccountBeta",
        "serviceAccountVcr",
        "projectGa",
        "projectBeta",
        "projectVcr",
        "projectNumberGa",
        "projectNumberBeta",
        "projectNumberVcr",
        "identityUserGa",
        "identityUserBeta",
        "identityUserVcr",
        "masterBillingAccountGa",
        "masterBillingAccountBeta",
        "masterBillingAccountVcr",
        "org2Ga",
        "org2Beta",
        "org2Vcr",
        "billingAccount",
        "billingAccount2",
        "custId",
        "org",
        "orgDomain",
        "region",
        "zone",
        "infraProject",
        "vcrBucketName",
        "credentialsGCS")
}

// getNestedProjectFromRoot allows you to retrieve a project that's 2 levels deep from the root project,
// Using the names of the parent and desired project.
// E.g. Root > Project A > Project B - you need to supply the inputs "Project A" and "Project B"
fun getNestedProjectFromRoot(root: Project, parentName: String, nestedProjectName: String): Project {
    // Find parent project within root
    val parent: Project = getSubProject(root, parentName)
    // Find subproject within parent identified above
    return getSubProject(parent, nestedProjectName)
}

// getSubProject allows you to retrieve a project nested inside a given project you've already identified,
// using its name.
fun getSubProject(parent: Project, subProjectName: String): Project {
    val subProject: Project? =  parent.subProjects.find { p->  p.name == subProjectName}
    if (subProject == null) {
        fail("Could not find the $subProjectName project inside ${parent.name}")
    }
    return subProject!!
}

// getBuildFromProject allows you to retrieve a build configuration from an identified project using its name
fun getBuildFromProject(parentProject: Project, buildName: String): BuildType {
    val buildType: BuildType?  = parentProject.buildTypes.find { p->  p.name == buildName}
    if (buildType == null) {
        fail("Could not find the '$buildName' build in project ${parentProject.name}")
    }
    return buildType!!
}