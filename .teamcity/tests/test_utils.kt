/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package tests

import builds.AllContextParameters

const val gaProjectName = "Google"
const val betaProjectName = "Google Beta"
const val nightlyTestsProjectName = "Nightly Tests"
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