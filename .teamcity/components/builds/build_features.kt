// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import jetbrains.buildServer.configs.kotlin.BuildFeatures
import jetbrains.buildServer.configs.kotlin.buildFeatures.GolangFeature

// NOTE: this file includes Extensions of the Kotlin DSL class BuildFeature
// This allows us to reuse code in the config easily, while ensuring the same build features can be used across builds.
// See the class's documentation: https://teamcity.jetbrains.com/app/dsl-documentation/root/build-feature/index.html


const val UseTeamCityGoTest = false

fun BuildFeatures.golang() {
    if (UseTeamCityGoTest) {
        feature(GolangFeature {
            testFormat = "json"
        })
    }
}