/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.AllContextParameters
import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.DslContext
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.SharedResources
import jetbrains.buildServer.configs.kotlin.triggers.FinishBuildTrigger
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Assert.fail

const val gaProjectName = "Google"
const val betaProjectName = "Google Beta"
const val nightlyTestsProjectName = "Nightly Tests"
const val weeklyDiffTestsProjectName = "Weekly Diff Tests"
const val mmUpstreamProjectName = "Upstream MM Testing"
const val globalSweepersProjectName = "Global Sweepers"

fun testContextParameters(projectId: String = "TeamCityTests"): AllContextParameters {
    DslContext.projectId = AbsoluteId(projectId)
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
        "chronicleInstanceIdGa",
        "chronicleInstanceIdBeta",
        "chronicleInstanceIdVcr",
        "vmwareengineProjectGa",
        "vmwareengineProjectBeta",
        "vmwareengineProjectVcr",
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

// getNestedProjectFromRoot allows you to retrieve a project that's 2 levels deep from the root project,
// Using the names of the parent and desired project.
// E.g. Root > Project A > Project B - you need to supply the inputs "Project A" and "Project B"
fun getNestedProjectFromRoot(root: Project, parentName: String, nestedProjectName: String): Project {
    // Find parent project within root
    val parent: Project = getSubProject(root, parentName)
    // Find subproject within parent identified above
    return getSubProject(parent, nestedProjectName)
}

// getSubProject allows you to retrieve a project nested inside a given project you've already identified,
// using its name.
fun getSubProject(parent: Project, subProjectName: String): Project {
    val subProject: Project? =  parent.subProjects.find { p->  p.name == subProjectName}
    if (subProject == null) {
        fail("Could not find the $subProjectName project inside ${parent.name}")
    }
    return subProject!!
}

// getBuildFromProject allows you to retrieve a build configuration from an identified project using its name
fun getBuildFromProject(parentProject: Project, buildName: String): BuildType {
    val buildType: BuildType?  = parentProject.buildTypes.find { p->  p.name == buildName}
    if (buildType == null) {
        fail("Could not find the '$buildName' build in project ${parentProject.name}")
    }
    return buildType!!
}

fun assertFinishTrigger(build: BuildType, source: BuildType, branch: String) {
    assertEquals("${build.name} should have exactly one trigger", 1, build.triggers.items.size)
    val trigger = build.triggers.items.single()
    assertTrue("${build.name} should use a finish-build trigger", trigger is FinishBuildTrigger)
    trigger as FinishBuildTrigger
    assertEquals("${build.name} should watch the source build's resolved ID", source.id!!.value, trigger.buildType)
    assertEquals("${build.name} should watch the configured branch", "+:$branch", trigger.branchFilter)
    assertEquals("${build.name} should not require a successful source build", false, trigger.successfulOnly)
}

fun assertSnapshotDependencies(build: BuildType, sources: List<BuildType>, failureAction: FailureAction) {
    val dependencies = build.dependencies.items
    assertEquals("${build.name} should have exactly the expected dependencies", sources.size, dependencies.size)
    assertEquals(
        "${build.name} should depend on the registered source builds",
        sources.map { it.id!!.value }.toSet(),
        dependencies.map { it.buildTypeId.id!!.value }.toSet()
    )
    dependencies.forEach { dependency ->
        val snapshot = requireNotNull(dependency.snapshot) { "${build.name} should use snapshot dependencies" }
        assertEquals("${build.name} should preserve dependency failure handling", failureAction, snapshot.onDependencyFailure)
        assertEquals("${build.name} should preserve dependency cancellation handling", failureAction, snapshot.onDependencyCancel)
        assertTrue("${build.name} should synchronize source revisions", snapshot.synchronizeRevisions)
    }
}

fun assertSharedResourceLocks(build: BuildType, expected: SharedResources) {
    val features = build.features.items.filterIsInstance<SharedResources>()
    assertEquals("${build.name} should have one shared-resource feature", 1, features.size)
    assertEquals(
        "${build.name} should acquire the expected locks",
        expected.params.single { it.name == "locks-param" }.value,
        features.single().params.single { it.name == "locks-param" }.value
    )
}