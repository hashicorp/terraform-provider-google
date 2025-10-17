/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects

import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.AllContextParameters
import builds.readOnlySettings
import generated.PackagesListBeta
import generated.PackagesListGa
import generated.ServicesListBeta
import generated.ServicesListGa
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.sharedResource
import projects.feature_branches.featureBranchResourceIdentitySubProject

// googleCloudRootProject returns a root project that contains a subprojects for the GA and Beta version of the
// Google provider. There are also resources to help manage the test projects used for acceptance tests.
fun googleCloudRootProject(allConfig: AllContextParameters): Project {

    return Project{

        description = "Contains all testing projects for the GA and Beta versions of the Google provider."

        // Registering the VCS roots used by subprojects
        vcsRoot(vcs_roots.HashiCorpVCSRootGa)
        vcsRoot(vcs_roots.HashiCorpVCSRootBeta)
        vcsRoot(vcs_roots.ModularMagicianVCSRootGa)
        vcsRoot(vcs_roots.ModularMagicianVCSRootBeta)

        features {
            // For controlling sweeping of the GA nightly test project
            sharedResource {
                id = "GA_NIGHTLY_SERVICE_LOCK_SHARED_RESOURCE"
                name = SharedResourceNameGa
                enabled = true
                resourceType = customValues(getPackageNameList(ServicesListGa) + getPackageNameList(PackagesListGa))
            }
            // For controlling sweeping of the Beta nightly test project
            sharedResource {
                id = "BETA_NIGHTLY_SERVICE_LOCK_SHARED_RESOURCE"
                name = SharedResourceNameBeta
                enabled = true
                resourceType = customValues(getPackageNameList(ServicesListBeta) + getPackageNameList(PackagesListBeta))
            }
            // For controlling sweeping of the PR testing project
            sharedResource {
                id = "PR_SERVICE_LOCK_SHARED_RESOURCE"
                name = SharedResourceNameVcr
                enabled = true
                resourceType = customValues(getPackageNameList(ServicesListBeta) + getPackageNameList(PackagesListBeta)) // Use Beta list of services here, assuming Beta is a superset of GA
            }
        }

        // Projects required for nightly testing, testing MM upstreams, and sweepers
        subProject(googleSubProjectGa(allConfig))
        subProject(googleSubProjectBeta(allConfig))
        subProject(projectSweeperSubProject(allConfig))
        subProject(featureBranchResourceIdentitySubProject(allConfig))

        // Feature branch-testing projects - these will be added and removed as needed

        params {
            readOnlySettings()
        }
    }
}

fun getPackageNameList(servicesList: Map<String, Map<String,String>>): List<String> {
    var serviceNameList: ArrayList<String> = arrayListOf()
    servicesList.forEach{ s ->
        var serviceName = s.value.getValue("name").toString()
        serviceNameList.add(serviceName)
    }
    return serviceNameList
}
