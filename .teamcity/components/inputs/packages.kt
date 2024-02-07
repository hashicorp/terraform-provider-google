// This file is controlled by MMv1, any changes made here will be overwritten

package generated

var PackagesList = mapOf(
    "envvar" to mapOf(
        "name" to "envvar",
        "displayName" to "Environment Variables",
        "path" to "./google/envvar"
    ),
    "fwmodels" to mapOf(
        "name" to "fwmodels",
        "displayName" to "Framework Models",
        "path" to "./google/fwmodels"
    ),
    "fwprovider" to mapOf(
        "name" to "fwprovider",
        "displayName" to "Framework Provider",
        "path" to "./google/fwprovider"
    ),
    "fwresource" to mapOf(
        "name" to "fwresource",
        "displayName" to "Framework Resource",
        "path" to "./google/fwresource"
    ),
    "fwtransport" to mapOf(
        "name" to "fwtransport",
        "displayName" to "Framework Transport",
        "path" to "./google/fwtransport"
    ),
    "provider" to mapOf(
        "name" to "provider",
        "displayName" to "SDK Provider",
        "path" to "./google/provider"
    ),
    "transport" to mapOf(
        "name" to "transport",
        "displayName" to "Transport",
        "path" to "./google/transport"
    ),
    "google" to mapOf(
        "name" to "google",
        "displayName" to "Google",
        "path" to "./google"
    )
)

var SweepersList = mapOf(
    "sweeper" to mapOf(
        "name" to "sweeper",
        "displayName" to "Sweeper",
        "path" to "./google/sweeper"
    )
)

fun GetPackageNameList(): List<String> {
    var packageNameList: ArrayList<String> = arrayListOf()
    PackagesList.forEach{ p ->
        var packageName = p.value.getValue("name").toString()
        packageNameList.add(packageName)
    }
    return packageNameList
}
