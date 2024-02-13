/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

fun replaceCharsId(id: String): String{
    var newId = id.replace("-", "")
    newId = newId.replace(" ", "_")
    newId = newId.uppercase()

    return newId
}