/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.ParametrizedWithType

class ClientConfiguration(var custId: String,
                          var org: String,
                          val org2 : String,
                          val billingAccount : String,
                          val billingAccount2: String,
                          val masterBillingAccount: String,
                          val credentials : String,
                          val project : String,
                          val orgDomain : String,
                          val projectNumber : String,
                          val region : String,
                          val serviceAccount : String,
                          val zone : String,
                          val firestoreProject : String,
                          val identityUser : String ) {
}

// ParametrizedWithType.ConfigureGoogleSpecificTestParameters allows build configs to be created
// with the environment variables needed to configure the provider and/or configure test code.
// Extension of ParametrizedWithType. For docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/parametrized-with-type/index.html
fun ParametrizedWithType.ConfigureGoogleSpecificTestParameters(config: ClientConfiguration) {
    hiddenPasswordVariable("env.GOOGLE_CUST_ID", config.custId, "The ID of the Google Identity Customer")
    hiddenPasswordVariable("env.GOOGLE_ORG", config.org, "The Google Organization Id")
    hiddenPasswordVariable("env.GOOGLE_ORG_2", config.org2, "The second Google Organization Id")
    hiddenPasswordVariable("env.GOOGLE_BILLING_ACCOUNT", config.billingAccount, "The billing account associated with the first google organization")
    hiddenPasswordVariable("env.GOOGLE_BILLING_ACCOUNT_2", config.billingAccount2, "The billing account associated with the second google organization")
    hiddenPasswordVariable("env.GOOGLE_MASTER_BILLING_ACCOUNT", config.masterBillingAccount, "The master billing account")
    hiddenVariable("env.GOOGLE_PROJECT", config.project, "The google project for this build")
    hiddenVariable("env.GOOGLE_ORG_DOMAIN", config.orgDomain, "The org domain")
    hiddenVariable("env.GOOGLE_PROJECT_NUMBER", config.projectNumber, "The project number associated with the project")
    hiddenVariable("env.GOOGLE_REGION", config.region, "The google region to use")
    hiddenVariable("env.GOOGLE_SERVICE_ACCOUNT", config.serviceAccount, "The service account")
    hiddenVariable("env.GOOGLE_ZONE", config.zone, "The google zone to use")
    hiddenVariable("env.GOOGLE_FIRESTORE_PROJECT", config.firestoreProject, "The project to use for firestore")
    hiddenVariable("env.GOOGLE_IDENTITY_USER", config.identityUser, "The user for the identity platform")
    hiddenPasswordVariable("env.GOOGLE_CREDENTIALS", config.credentials, "The Google credentials for this test runner")
}