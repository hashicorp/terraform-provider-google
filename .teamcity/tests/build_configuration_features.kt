/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.UseTeamCityGoTest
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert.assertTrue
import org.junit.Assert.fail
import org.junit.Test
import projects.googleCloudRootProject

class BuildConfigurationFeatureTests {
    @Test
    fun buildShouldFailOnError() {
        val project = googleCloudRootProject(testContextParameters())
        // Find Google (GA) project
        var gaProject: Project? =  project.subProjects.find { p->  p.name == gaProjectName}
        if (gaProject == null) {
            fail("Could not find the Google (GA) project")
        }
        // Find Google Beta project
        var betaProject: Project? =  project.subProjects.find { p->  p.name == betaProjectName}
        if (betaProject == null) {
            fail("Could not find the Google (GA) project")
        }

        (gaProject!!.subProjects + betaProject!!.subProjects).forEach{p ->
            p.buildTypes.forEach{bt ->
                assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
            }
        }
    }

    @Test
    fun buildShouldHaveGoTestFeature() {
        val project = googleCloudRootProject(testContextParameters())
        // Find Google (GA) project
        var gaProject: Project? =  project.subProjects.find { p->  p.name == gaProjectName}
        if (gaProject == null) {
            fail("Could not find the Google (GA) project")
        }
        // Find Google Beta project
        var betaProject: Project? =  project.subProjects.find { p->  p.name == betaProjectName}
        if (betaProject == null) {
            fail("Could not find the Google (GA) project")
        }

        (gaProject!!.subProjects + betaProject!!.subProjects).forEach{p ->
            var exists: ArrayList<Boolean> = arrayListOf()
            p.buildTypes.forEach{bt ->
                bt.features.items.forEach { f ->
                    exists.add(f.type == "golang")
                }
            }
            if (exists.contains(false) && UseTeamCityGoTest){
                fail("Project ${p.name} contains build configurations that don't use the Go Test feature")
            }
        }
    }

    @Test
    fun nonVCRBuildShouldHaveSaveArtifactsToGCS() {
        val project = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        var gaNightlyTestProject = getSubProject(project, gaProjectName, nightlyTestsProjectName)

        // Find GA MM Upstream project
        var gaMMUpstreamProject = getSubProject(project, gaProjectName, mmUpstreamProjectName)

        // Find Beta nightly test project
        var betaNightlyTestProject = getSubProject(project, betaProjectName, nightlyTestsProjectName)

        // Find Beta MM Upstream project
        var betaMMUpstreamProject = getSubProject(project, betaProjectName, mmUpstreamProjectName)

        (gaNightlyTestProject.buildTypes + gaMMUpstreamProject.buildTypes + betaNightlyTestProject.buildTypes + betaMMUpstreamProject.buildTypes).forEach{bt ->
            var found: Boolean = false
            for (item in bt.steps.items) {
                if (item.name == "Tasks after running nightly tests: push artifacts(debug logs) to GCS") {
                    found = true
                    break
                }
            }
            // service sweeper does not contain push artifacts to GCS step
            if (bt.name != "Service Sweeper") {
                assertTrue("Build configuration `${bt.name}` contains a build step that pushes artifacts to GCS", found)
            }
        }
    }
}
