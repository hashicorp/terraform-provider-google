/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.AllContextParameters
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert

const val gaProjectName = "Google"
const val betaProjectName = "Google Beta"
const val nightlyTestsProjectName = "Nightly Tests"
const val mmUpstreamProjectName = "MM Upstream Testing"
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
        "firestoreProjectGa",
        "firestoreProjectBeta",
        "firestoreProjectVcr",
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
        "vcrBucketName")
}

fun getSubProject(rootProject: Project, parentProjectName: String, subProjectName: String): Project {
    // Find parent project within root
    var parentProject: Project? =  rootProject.subProjects.find { p->  p.name == parentProjectName}
    if (parentProject == null) {
        Assert.fail("Could not find the $parentProjectName project")
    }
    // Find subproject within parent identified above
    var subProject: Project?  = parentProject!!.subProjects.find { p->  p.name == subProjectName}
    if (subProject == null) {
        Assert.fail("Could not find the $subProjectName project")
    }

    return subProject!!
}