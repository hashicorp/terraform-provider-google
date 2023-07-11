// this file is copied from mmv1, any changes made here will be overwritten

package tests

import ClientConfiguration

fun TestConfiguration() : ClientConfiguration {
    return ClientConfiguration("custId", "org", "org2", "billingAccount", "billingAccount2", "masterBillingAccount", "credentials", "project", "orgDomain", "projectNumber", "region", "serviceAccount", "zone", "firestoreProject", "identityUser")
}