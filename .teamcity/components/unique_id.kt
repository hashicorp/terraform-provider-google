/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

fun replaceCharsId(id: String): String{
    var newId = id.replace("-", "")
    newId = newId.replace(" ", "_")
    newId = newId.uppercase()

    return newId
}