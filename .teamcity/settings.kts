/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

import projects.googleCloudRootProject
import builds.AllContextParameters
import jetbrains.buildServer.configs.kotlin.*

version = "2023.11"

// The code below pulls context parameters from the TeamCity project.
// Context parameters aren't stored in VCS, and are managed manually.
// Due to this, the code needs to explicitly pull in values via the DSL and pass the values into other code.
// For DslContext docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/dsl-context/index.html

// Context parameters below are used to set ENVs needed for acceptance tests within the build configurations.

// GOOGLE_CREDENTIALS
val credentialsGa   = DslContext.getParameter("credentialsGa", "")
val credentialsBeta = DslContext.getParameter("credentialsBeta", "")
val credentialsVcr  = DslContext.getParameter("credentialsVcr", "")
// GOOGLE_SERVICE_ACCOUNT
val serviceAccountGa   = DslContext.getParameter("serviceAccountGa", "")
val serviceAccountBeta = DslContext.getParameter("serviceAccountBeta", "")
val serviceAccountVcr  = DslContext.getParameter("serviceAccountVcr", "")
// GOOGLE_PROJECT & GOOGLE_PROJECT_NUMBER
val projectGa         = DslContext.getParameter("projectGa", "")
val projectBeta       = DslContext.getParameter("projectBeta", "")
val projectVcr        = DslContext.getParameter("projectVcr", "")
val projectNumberGa   = DslContext.getParameter("projectNumberGa", "")
val projectNumberBeta = DslContext.getParameter("projectNumberBeta", "")
val projectNumberVcr  = DslContext.getParameter("projectNumberVcr", "")
// GOOGLE_IDENTITY_USER
val identityUserGa   = DslContext.getParameter("identityUserGa", "")
val identityUserBeta = DslContext.getParameter("identityUserBeta", "")
val identityUserVcr  = DslContext.getParameter("identityUserVcr", "")
// GOOGLE_FIRESTORE_PROJECT
val firestoreProjectGa   = DslContext.getParameter("firestoreProjectGa", "")
val firestoreProjectBeta = DslContext.getParameter("firestoreProjectBeta", "")
val firestoreProjectVcr  = DslContext.getParameter("firestoreProjectVcr", "")
// GOOGLE_MASTER_BILLING_ACCOUNT
val masterBillingAccountGa   = DslContext.getParameter("masterBillingAccountGa", "")
val masterBillingAccountBeta = DslContext.getParameter("masterBillingAccountBeta", "")
val masterBillingAccountVcr  = DslContext.getParameter("masterBillingAccountVcr", "")
// GOOGLE_ORG_2
val org2Ga   = DslContext.getParameter("org2Ga", "")
val org2Beta = DslContext.getParameter("org2Beta", "")
val org2Vcr  = DslContext.getParameter("org2Vcr", "")

// Values that are the same across GA, Beta, and VCR testing environments
val billingAccount  = DslContext.getParameter("billingAccount", "")   // GOOGLE_BILLING_ACCOUNT
val billingAccount2 = DslContext.getParameter("billingAccount2", "")  // GOOGLE_BILLING_ACCOUNT_2
val custId          = DslContext.getParameter("custId", "")           // GOOGLE_CUST_ID
val org             = DslContext.getParameter("org", "")              // GOOGLE_ORG
val orgDomain       = DslContext.getParameter("orgDomain", "")        // GOOGLE_ORG_DOMAIN
val region          = DslContext.getParameter("region", "")           // GOOGLE_REGION
val zone            = DslContext.getParameter("zone", "")             // GOOGLE_ZONE

// Used for recording VCR cassettes
val infraProject             = DslContext.getParameter("infraProject", "") // GOOGLE_INFRA_PROJECT
val vcrBucketName            = DslContext.getParameter("vcrBucketName", "") // VCR_BUCKET_NAME


var allContextParams = AllContextParameters(
    credentialsGa,
    credentialsBeta,
    credentialsVcr,
    serviceAccountGa,
    serviceAccountBeta,
    serviceAccountVcr,
    projectGa,
    projectBeta,
    projectVcr,
    projectNumberGa,
    projectNumberBeta,
    projectNumberVcr,
    identityUserGa,
    identityUserBeta,
    identityUserVcr,
    firestoreProjectGa,
    firestoreProjectBeta,
    firestoreProjectVcr,
    masterBillingAccountGa,
    masterBillingAccountBeta,
    masterBillingAccountVcr,
    org2Ga,
    org2Beta,
    org2Vcr,
    billingAccount,
    billingAccount2,
    custId,
    org,
    orgDomain,
    region,
    zone,
    infraProject,
    vcrBucketName
)

// This is the entry point of the code in .teamcity/
// See https://teamcity.jetbrains.com/app/dsl-documentation/root/project.html
project(googleCloudRootProject(allContextParams))