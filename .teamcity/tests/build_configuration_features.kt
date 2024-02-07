// This file is controlled by MMv1, any changes made here will be overwritten

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
}
