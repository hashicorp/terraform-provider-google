# Using Configuration Files with TeamCity

The testing environment in TeamCity can be found at https://hashicorp.teamcity.com/

Contents:
* [Using configuration files for the first time in a project](#using-configuration-files-for-the-first-time-in-a-project)
* [Updating configuration files](#updating-configuration-files)
* [Editing configuration files](#editing-configuration-files)
* [Pushing configuration changes to TeamCity](#pushing-configuration-changes-to-teamcity)


## Using configuration files for the first time in a project

At the start, we need to create a project in TeamCity and connect it to our hashicorp/terraform-provider-google repository, so TeamCity can pull in the configuration files to set that project's settings. This will cause all the resources defined in the configuration files to be created within that parent project.

The hashicorp/terraform-provider-google repository contains the entirety of the configuration required for testing both GA and Beta versions of the provider. The process described below only needs to be completed once to provision all required projects.

**Note: The steps below require the user to have Admin permissions in TeamCity.**

### 1) Create a parent project to be populated using the TeamCity configuration in GitHub
 
* When viewing the [Terraform Providers project in TeamCity](https://hashicorp.teamcity.com/project/TerraformProviders), click the dropdown next to `Edit project...` in the top right, and select `New subproject...`
* On the `Create Project` page:
* Select `From a repository URL`.
* Ensure Parent project is 'Terraform Providers'.
* Set Repository URL to `https://github.com/hashicorp/terraform-provider-google.git`.
* Leave Username and Password / access token blank; the repository is public access.
* Click `Proceed`. You'll then be taken to the `Create Project From URL` page
* On the next page, select `Import settings from .teamcity/settings.kts and enable synchronization with the VCS repository`
* Set `Project name` to an appropriate name, e.g. "Google Cloud"
* Ensure `Default branch` is `refs/heads/main`.
* Empty the `Branch specification` text box; we don't want to monitor branches other than `main`.
* Click `Proceed`.
    * The button will become disabled and a progress spinner will appear.
    * Wait for the page to reload after the settings are successfully loaded into your new project.

What we've done so far: The process above creates a new project and a [VCS Root](https://www.jetbrains.com/help/teamcity/configuring-vcs-roots.html), through which TeamCity can pull configuration files from GitHub and apply those files to set the new project's settings. The project's settings will be synchronised with the configuration in GitHub, so if you updated the configuration files in GitHub those changes would automatically be reflected in TeamCity.

The next step is provide some input values that the configuration needs to fully function. 


### 2) Provide inputs via Context Parameters and Tokens

* Following the previous steps you should have been redirected to the Project settings page for your new project. If needed, navigate to this page manually.
* Next, click `Versioned Settings` in the left menu
    * You should see that, from the previous steps, versioned settings is active and `Synchronization enabled` is selected.
* Click on the `Context Parameters` tab at the top of the page
    * On this page you can add key:value pairs, which are used as input to the configuration in `.teamcity/settings.kts`
* First, we need to enter credentials information in a special way that keeps the values secure. Here is the process for the credentials to the GA nightly test project:
    * Find the credentials JSON file for that project, ensure that the value has no newlines
    * Click the dropdown menu next to `Actions` in the top right and click `Generate token for a secure value...`
    * Paste the credentials JSON string into the `Secure value` field nd click `Generate Token`.
    * Copy the value shown below the field, which should look like `credentialsJSON:<uuid>`
    * Click `Close`
    * Still on the `Context Parameters` page, create a new entry where the key is `credentialsGa` and the value is that secure token
* Repeat the above process for all 3 GCP projects, one for each row below:

<br/>

| Name | Notes |
|---|---|
| credentialsGa | USE SECURE TOKENS TO ENTER THIS VALUE. Used to set the [GOOGLE_CREDENTIALS](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L15) environment variable in acceptance tests - GA nightly test project specific |
| credentialsBeta | USE SECURE TOKENS TO ENTER THIS VALUE. Used to set the [GOOGLE_CREDENTIALS](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L15) environment variable in acceptance tests - Beta nightly test project specific |
| credentialsVcr | USE SECURE TOKENS TO ENTER THIS VALUE. Used to set the [GOOGLE_CREDENTIALS](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L15) environment variable in acceptance tests - VCR test project specific |


<br/>

* Next, enter values for all the context parameters below. These don't use secure tokens:

<br/>


| Name | Notes |
|---|---|
| serviceAccountGa | Used to set the [GOOGLE_SERVICE_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L69-L71) environment variable in acceptance tests - GA specific |
| serviceAccountBeta | Used to set the [GOOGLE_SERVICE_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L69-L71) environment variable in acceptance tests - Beta specific |
| serviceAccountVcr | Used to set the [GOOGLE_SERVICE_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L69-L71) environment variable in acceptance tests - VCR specific |
| projectGa | Used to set the [GOOGLE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L27)  environment variable in acceptance tests- GA specific |
| projectBeta | Used to set the [GOOGLE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L27)  environment variable in acceptance tests- Beta specific |
| projectVcr | Used to set the [GOOGLE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L27)  environment variable in acceptance tests- VCR specific |
| projectNumberGa | Used to set the [GOOGLE_PROJECT_NUMBER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L22-L24) environment variable in acceptance tests - GA specific |
| projectNumberBeta | Used to set the [GOOGLE_PROJECT_NUMBER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L22-L24) environment variable in acceptance tests - Beta specific |
| projectNumberVcr | Used to set the [GOOGLE_PROJECT_NUMBER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L22-L24) environment variable in acceptance tests - VCR specific |
| identityUserGa | Used to set the [GOOGLE_IDENTITY_USER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L58-L63) environment variable in acceptance tests - GA specific |
| identityUserBeta | Used to set the [GOOGLE_IDENTITY_USER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L58-L63) environment variable in acceptance tests - Beta specific |
| identityUserVcr | Used to set the [GOOGLE_IDENTITY_USER](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L58-L63) environment variable in acceptance tests - VCR specific |
| firestoreProjectGa | Used to set the [GOOGLE_FIRESTORE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L32-L34) environment variable in acceptance tests - GA specific |
| firestoreProjectBeta | Used to set the [GOOGLE_FIRESTORE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L32-L34) environment variable in acceptance tests - Beta specific |
| firestoreProjectVcr | Used to set the [GOOGLE_FIRESTORE_PROJECT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L32-L34) environment variable in acceptance tests - VCR specific |
| masterBillingAccountGa | Used to set the [GOOGLE_MASTER_BILLING_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L87-L91) environment variable in acceptance tests - GA specific |
| masterBillingAccountBeta | Used to set the [GOOGLE_MASTER_BILLING_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L87-L91) environment variable in acceptance tests - Beta specific |
| masterBillingAccountVcr | Used to set the [GOOGLE_MASTER_BILLING_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L87-L91) environment variable in acceptance tests - VCR specific |
| org2Ga | Used to set the [GOOGLE_ORG_2](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L73-L75) environment variable in acceptance tests - GA specific |
| org2Beta | Used to set the [GOOGLE_ORG_2](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L73-L75) environment variable in acceptance tests - Beta specific |
| org2Vcr | Used to set the [GOOGLE_ORG_2](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L73-L75) environment variable in acceptance tests - VCR specific |
| billingAccount | Used to set the [GOOGLE_BILLING_ACCOUNT](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L81-L85) ALL environment variable in acceptance tests |
| billingAccount2 | Used to set the [GOOGLE_BILLING_ACCOUNT_2](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_test.go#L78-L79) environment variable in ALL acceptance tests |
| custId | Used to set the [GOOGLE_CUST_ID](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L52-L56) environment variable in ALL acceptance tests |
| org | Used to set the [GOOGLE_ORG](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L48-L50) environment variable in ALL acceptance tests |
| orgDomain | Used to set the [GOOGLE_ORG_DOMAIN](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L65-L67) environment variable in ALL acceptance tests |
| region | Used to set the [GOOGLE_REGION](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L37)  environment variable in ALL acceptance tests|
| zone | Used to set the [GOOGLE_ZONE](https://github.com/GoogleCloudPlatform/magic-modules/blob/94a3f91d75ee823c521a0d8d3984a1493fa0926a/mmv1/third_party/terraform/envvar/envvar_utils.go#L43)  environment variable in ALL acceptance tests|
| infraProject | Used to set an environment variable `GOOGLE_INFRA_PROJECT` that's used to control the bucket used by VCR testing builds |
| vcrBucketName | Used to set an environment variable `VCR_BUCKET_NAME` that's used to control the bucket used by VCR testing builds |


* Click `Save`
    * TeamCity will take some time to process the new values. The status of loading the new values will show at the bottom of the Versioned Settings page
    * You can check if values are being set by navigating to an acceptance test or VCR recording build configuration clicking `View configuration settings` in the top right, and clicking `Parameters` in the left menu. Credentials values will be obscured as `******` but all other values will be human-readable.


---

## Editing configuration files

See [CONTRIBUTION_GUIDE.md](./CONTRIBUTION_GUIDE.md)

---

## Pushing configuration changes to TeamCity

When changes to files in `.teamcity/` are merged into the `main` branch of the hashicorp/terraform-provider-google repository TeamCity will detect these changes and attempt to apply them.

### Speeding up new changes being reflected in TeamCity

* Navigate to the `Versioned Settings` page on the parent project.
* Click the `Change Log` tab at the top.
* Click the `Check for changes` button.
    * This will cause TeamCity to immediately look for changes in the configuration, ahead of the normal schedule.

### Bugs in the TeamCity configuration files

If an update to the TeamCity configuration isn't valid, for example there is a syntax issue in a Kotlin file, then TeamCity will detect the problem and revert to using the config in the last known good commit.

In scenarios liek this, you will see an error shown at the bottom of the `Versioned Settings` page. Builds will continue to run using the last version of the configuration files.

To avoid bugs reaching TeamCity, validate all changes to the TeamCity configuration files locally before merging a PR:
* Generate the GA provider from Magic Modules
* In the GA provider repo, run `cd .teamcity`
* Validate the configuration using `make validate`