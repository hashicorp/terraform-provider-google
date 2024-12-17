/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
package components.projects.feature_branches
 
import generated.ServicesListGa
import generated.ServicesListBeta

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

// This function is used to get the services list for a given version. Typically used in feature branch builds for testing very specific services only.
fun getServicesList(Services: Array<String>, version: String): Map<String,Map<String,String>> {
    if (Services.isEmpty()) {
        throw Exception("No services found for version $version")
    }

    var servicesList = mutableMapOf<String,Map<String,String>>()
    for (service in Services) {
        if (version == "GA" || version == "GA-MM") {
            servicesList[service] = ServicesListGa.getOrElse(service) { throw Exception("Service $service not found") }
        } else if (version == "Beta" || version == "Beta-MM") {
            servicesList[service] = ServicesListBeta.getOrElse(service) { throw Exception("Service $service not found") }
        } else {
            throw Exception("Invalid version $version")
        }
    }

    when (version) {
        "GA" -> servicesList
        "Beta" -> {
            servicesList.mapValues { (_, value) ->
                value + mapOf(
                        "displayName" to "${value["displayName"]} - Beta"
                )
            }.toMutableMap()
        }
        "GA-MM" -> {
            servicesList.mapValues { (_, value) ->
                value + mapOf(
                        "displayName" to "${value["displayName"]} - MM"
                )
            }.toMutableMap()
        }
        "Beta-MM" -> {
            servicesList.mapValues { (_, value) ->
                value + mapOf(
                        "displayName" to "${value["displayName"]} - Beta - MM"
                )
            }.toMutableMap()
        }
        else -> throw Exception("Invalid version $version")
    }.also { servicesList = it as MutableMap<String, Map<String, String>> }

    return servicesList
}