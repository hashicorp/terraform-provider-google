/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

import Google
import ClientConfiguration
import jetbrains.buildServer.configs.kotlin.*

version = "2023.05"

// The code below pulls context parameters from the TeamCity project.
// Context parameters aren't stored in VCS, and are managed manually.
// Due to this, the code needs to explicitly pull in values via the DSL and pass the values into other code.
// For DslContext docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/dsl-context/index.html

// Values of these context parameters are used to set ENVs needed for acceptance tests within the build configurations.
var custId = DslContext.getParameter("custId", "")
var org = DslContext.getParameter("org", "")
var org2 = DslContext.getParameter("org2", "")
var billingAccount = DslContext.getParameter("billingAccount", "")
var billingAccount2 = DslContext.getParameter("billingAccount2", "")
var masterBillingAccount = DslContext.getParameter("masterBillingAccount", "")
var project = DslContext.getParameter("project", "")
var orgDomain = DslContext.getParameter("orgDomain", "")
var projectNumber = DslContext.getParameter("projectNumber", "")
var region = DslContext.getParameter("region", "")
var serviceAccount = DslContext.getParameter("serviceAccount", "")
var zone = DslContext.getParameter("zone", "")
var credentials = DslContext.getParameter("credentials", "")
var firestoreProject = DslContext.getParameter("firestoreProject", "")
var identityUser = DslContext.getParameter("identityUser", "")

// Get details of the VCS Root that's manually made when VCS is first
// connected to the Project in TeamCity
var manualVcsRoot = DslContext.settingsRootId

// Values of these context parameters change configuration code behaviour.
var environment = DslContext.getParameter("environment", "default")
var branchRef = DslContext.getParameter("branch", "refs/heads/main")
var projDescription = DslContext.getParameter("description", "")

var clientConfig = ClientConfiguration(custId, org, org2, billingAccount, billingAccount2, masterBillingAccount, credentials, project, orgDomain, projectNumber, region, serviceAccount, zone, firestoreProject, identityUser)

// This is the entry point of the code in .teamcity/
// See https://teamcity.jetbrains.com/app/dsl-documentation/root/project.html
project(Google(environment, projDescription, manualVcsRoot, branchRef, clientConfig))
