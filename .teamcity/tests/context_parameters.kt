/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import builds.AccTestConfiguration
import builds.getBetaAcceptanceTestConfig
import builds.getGaAcceptanceTestConfig
import builds.getVcrAcceptanceTestConfig
import org.junit.Assert.fail
import org.junit.Test
import kotlin.reflect.full.memberProperties

class ContextParameterHandlingTests {
    @Test
    fun getGaAcceptanceTestConfig_returnsGaValuesOnly() {
        val config: AccTestConfiguration = getGaAcceptanceTestConfig(testContextParameters())
        for (prop in AccTestConfiguration::class.memberProperties) {
            val value = prop.get(config).toString()
            if (value.contains("Beta")||value.contains("Vcr")) {
                fail("Found config value $value which isn't a GA value")
            }
        }
    }

    @Test
    fun getBetaAcceptanceTestConfig_returnsBetaValuesOnly() {
        val config: AccTestConfiguration = getBetaAcceptanceTestConfig(testContextParameters())
        for (prop in AccTestConfiguration::class.memberProperties) {
            val value = prop.get(config).toString()
            if (value.contains("Ga")||value.contains("Vcr")) {
                fail("Found config value $value which isn't a Beta value")
            }
        }
    }

    @Test
    fun getVcrAcceptanceTestConfig_returnsVcrValuesOnly() {
        val config: AccTestConfiguration = getVcrAcceptanceTestConfig(testContextParameters())
        for (prop in AccTestConfiguration::class.memberProperties) {
            val value = prop.get(config).toString()
            if (value.contains("Ga")||value.contains("Beta")) {
                fail("Found config value $value which isn't a VCR value")
            }
        }
    }

}
