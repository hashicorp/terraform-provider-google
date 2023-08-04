/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

var packages = mapOf(
        "acctest" to "AccTest",
        "provider" to "SDK Provider",
        "fwprovider" to "Framework Plugin Provider",
        "services" to "Services",
        "tpgdclresource" to "TPG DCL Resource",
        "tpgiamresource" to "TPG IAM Resource",
        "tpgresource" to "TPG Resource",
        "transport" to "Transport",
        "verify" to "Verify"
)