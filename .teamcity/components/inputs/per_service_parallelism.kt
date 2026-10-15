/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package generated

// specifies the default level of parallelism for per-service-package test builds, unless overridden in ServiceParallelism below
const val DefaultServiceParallelism = 20

var ServiceParallelism = mapOf(
    // Packages with quota failures in nightly tests are kept at the previous default of 6.
    // gemini: Quota 'Maximum number of Code Repository Indexes' exceeded (limit 10)
    "gemini" to 6,
    // netapp: Quotas 'BackupVaultsPerRegion' (limit 10) and 'ReplicatedVolumesPerRegion' (limit 5) exceeded
    "netapp" to 5,
    // osconfig: OS Config API quota exhausted
    "osconfig" to 6,

    "looker" to 1,
    // dataprocmetastore: ;Create Bucket quota limit limit exceeded on the Google Cloud project. Please slow down execution and try again.
    "dataprocmetastore" to 6
)
