/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.AllContextParameters
import jetbrains.buildServer.BuildProject
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert

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

fun getSubProject(rootProject: Project, parentProjectName: String, subProjectName: String): Project {
    // Find parent project within root
    val parentProject: Project? =  rootProject.subProjects.find { p->  p.name == parentProjectName}
    if (parentProject == null) {
        Assert.fail("Could not find the $parentProjectName project")
    }
    // Find subproject within parent identified above
    val subProject: Project?  = parentProject!!.subProjects.find { p->  p.name == subProjectName}
    if (subProject == null) {
        Assert.fail("Could not find the $subProjectName project")
    }

    return subProject!!
}

fun getBuildFromProject(parentProject: Project, buildName: String): BuildType {
    val buildType: BuildType?  = parentProject!!.buildTypes.find { p->  p.name == buildName}
    if (buildType == null) {
        Assert.fail("Could not find the '$buildName' build in project ${parentProject.name}")
    }

    return buildType!!
}
