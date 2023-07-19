// this file is auto-generated with mmv1, any changes made here will be overwritten

import Google
import ClientConfiguration
import jetbrains.buildServer.configs.kotlin.*

version = "2023.05"

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
var environment = DslContext.getParameter("environment", "public")
var firestoreProject = DslContext.getParameter("firestoreProject", "")
var identityUser = DslContext.getParameter("identityUser", "")

// Get details of the VCS Root that's manually made when VCS is first
// connected to the Project in TeamCity
var manualVcsRoot = DslContext.settingsRootId

var clientConfig = ClientConfiguration(custId, org, org2, billingAccount, billingAccount2, masterBillingAccount, credentials, project, orgDomain, projectNumber, region, serviceAccount, zone, firestoreProject, identityUser)

project(Google(environment, manualVcsRoot, clientConfig))
