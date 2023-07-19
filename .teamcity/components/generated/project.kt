// this file is auto-generated with mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.AbsoluteId

const val providerName = "google"

fun Google(environment: String, manualVcsRoot: AbsoluteId, configuration: ClientConfiguration) : Project {
    return Project{

        var buildConfigs = buildConfigurationsForPackages(packages, providerName, "google", environment, manualVcsRoot, configuration)
        buildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
    }
}

fun buildConfigurationsForPackages(packages: Map<String, String>, providerName : String, path : String, environment: String, manualVcsRoot: AbsoluteId, config: ClientConfiguration): List<BuildType> {
    var list = ArrayList<BuildType>()

    packages.forEach { (packageName, displayName) ->
        if (packageName == "services") {
            var serviceList = buildConfigurationsForPackages(services, providerName, path+"/"+packageName, environment, manualVcsRoot, config)
            list.addAll(serviceList)
        } else {
            var defaultTestConfig = testConfiguration()

            var pkg = packageDetails(packageName, displayName, environment)
            var buildConfig = pkg.buildConfiguration(providerName, path, manualVcsRoot, true, defaultTestConfig.startHour, defaultTestConfig.parallelism, defaultTestConfig.daysOfWeek, defaultTestConfig.daysOfMonth)

            buildConfig.params.ConfigureGoogleSpecificTestParameters(config)

            list.add(buildConfig)
        }
    }

    return list
}

class testConfiguration(parallelism: Int = defaultParallelism, startHour: Int = defaultStartHour, daysOfWeek: String = defaultDaysOfWeek, daysOfMonth: String = defaultDaysOfMonth) {
    var parallelism = parallelism
    var startHour = startHour
    var daysOfWeek = daysOfWeek
    var daysOfMonth = daysOfMonth
}
