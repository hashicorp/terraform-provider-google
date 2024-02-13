/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import jetbrains.buildServer.configs.kotlin.BuildSteps
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep

fun BuildSteps.checkVcrEnvironmentVariables() {
    step(ScriptBuildStep {
        name = "Setup for running VCR tests: feedback about user-supplied environment variables and available CLI tools"
        scriptContent = """
            #!/bin/bash
            echo "VCR TESTING ENVIRONMENT CHECKS - ENVs and CLI TOOLS"
            if [ "${'$'}VCR_MODE" = "" ]; then
                echo "VCR_MODE is not set"
                exit 1
            fi
            if [ "${'$'}VCR_PATH" = "" ]; then
                echo "VCR_PATH is not set"
                exit 1
            fi
            if [ "${'$'}GOOGLE_INFRA_PROJECT" = "" ]; then
                echo "GOOGLE_INFRA_PROJECT is not set"
                exit 1
            fi
            if [ "${'$'}VCR_BUCKET_NAME" = "" ]; then
                echo "VCR_BUCKET_NAME is not set"
                exit 1
            fi
            if [ "${'$'}TEST" = "" ]; then
                echo "TEST is not set - set it to a value like ./google/..., ./google/services/... or ./google-beta/services/compute"
                exit 1
            fi
            if [ "${'$'}TESTARGS" = "" ]; then
                echo "TESTARGS is not set - set it to a value like -run=TestAccFoobar"
                exit 1
            fi

            if ! command -v gcloud &> /dev/null   
            then
                echo "gcloud CLI not found"
                exit 1
            fi

            if ! command -v gsutil &> /dev/null   
            then
                echo "gsutil CLI not found"
                exit 1
            fi
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}

fun BuildSteps.tagBuildToIndicateVcrMode() {
    step(ScriptBuildStep {
        name = "Set build tag to indicate if build is recording or replaying VCR cassettes"
        scriptContent = """
            #!/bin/bash
            echo "##teamcity[addBuildTag '${'$'}{VCR_MODE}']"
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}

fun BuildSteps.runVcrAcceptanceTests() {
    step(ScriptBuildStep {
        name = "Run Tests"
        scriptContent =  """
        echo "VCR Testing: Running acceptance tests"
        echo "TESTARGS = ${'$'}TESTARGS"
        echo "TEST = ${'$'}TEST"

        go test ${'$'}TEST -v ${'$'}TESTARGS -timeout="%TIMEOUT%h" -test.parallel="%PARALLELISM%" -ldflags="-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc"
        """.trimIndent()
    })
}

fun BuildSteps.runVcrTestRecordingSetup() {
    step(ScriptBuildStep {
        name = "Setup for running VCR tests: if in REPLAY mode, download existing cassettes"
        scriptContent = """
            #!/bin/bash
            echo "VCR Testing: Pre-test setup"
            echo "VCR_MODE: ${'$'}{VCR_MODE}"
            echo "VCR_PATH: ${'$'}{VCR_PATH}"
            
            # Ensure directory exists regardless of VCR mode
            mkdir -p ${'$'}VCR_PATH
            
            if [ "${'$'}VCR_MODE" = "RECORDING" ]; then
                echo "RECORDING MODE - skipping this build step; nothing needed from Cloud Storage bucket"
                exit 0
            fi

            echo "REPLAY MODE- retrieving cassettes from Cloud Storage bucket"

            # Authenticate gcloud CLI
            echo "${'$'}{GOOGLE_CREDENTIALS}" > google-account.json
            chmod 600 google-account.json
            gcloud auth activate-service-account --key-file=google-account.json

            # Pull files from GCS
            echo "Listing files present in gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/"
            gsutil ls -p ${'$'}GOOGLE_INFRA_PROJECT gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/

            echo "Copying files present in gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/"
            gsutil -m cp gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/* ${'$'}VCR_PATH

            # copy branch-specific cassettes over master. This might fail but that's ok if the folder doesnt exist
            export BRANCH_NAME=%teamcity.build.branch%
            echo "Attempting to copy branch-specific files, if they exist: gs://${'$'}VCR_BUCKET_NAME/beta/${'$'}BRANCH_NAME/fixtures/*"
            gsutil -m cp gs://${'$'}VCR_BUCKET_NAME/beta/${'$'}BRANCH_NAME/fixtures/* ${'$'}VCR_PATH

            echo "Listing files present in ${'$'}VCR_PATH:"
            ls ${'$'}VCR_PATH

            # Cleanup
            rm google-account.json
            gcloud auth application-default revoke
            gcloud auth revoke --all

            echo "Finished"
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}

fun BuildSteps.runVcrTestRecordingSaveCassettes() {
    step(ScriptBuildStep {
        name = "Tasks after running VCR tests: if in RECORDING mode, push new cassettes to GCS"
        scriptContent = """
            #!/bin/bash
            echo "VCR Testing: Post-test steps"
            echo "VCR_MODE: ${'$'}{VCR_MODE}"
            echo "VCR_PATH: ${'$'}{VCR_PATH}"

            if [ "${'$'}VCR_MODE" = "REPLAYING" ]; then
            echo "REPLAYING MODE - skipping this build step; nothing to be done"
            exit 0
            fi

            echo "RECORDING MODE - push new cassettes to Cloud Storage bucket"

            # Authenticate gcloud CLI
            echo "${'$'}{GOOGLE_CREDENTIALS}" > google-account.json
            chmod 600 google-account.json
            gcloud auth activate-service-account --key-file=google-account.json

            echo "Listing files present in ${'$'}VCR_PATH:"
            ls ${'$'}VCR_PATH

            # Fail early & explicitly if no files were created in the previous step
            export FILE_COUNT=${'$'}(ls ${'$'}VCR_PATH | wc -l)
            if test ${'$'}FILE_COUNT -le "0"; then
                echo "No cassette files found locally to upload to Cloud Storage bucket"
                exit 1
            fi

            export BRANCH_NAME=%teamcity.build.branch%
            if [ "${'$'}BRANCH_NAME" = "refs/heads/main" ]; then
                echo "Using main branch, so copying files to fixures/ in root of Cloud Storage bucket"

                echo "Listing files already present in gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/:"
                gsutil ls -p ${'$'}GOOGLE_INFRA_PROJECT gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/

                echo "Copying files to Cloud Storage bucket:"
                gsutil -m cp ${'$'}VCR_PATH/* gs://${'$'}VCR_BUCKET_NAME/beta/fixtures/
            else
                echo "Using ${'$'}BRANCH_NAME branch, so copying files to ${'$'}BRANCH_NAME/fixtures/ folder in Cloud Storage bucket"

                gsutil -m cp ${'$'}VCR_PATH/* gs://${'$'}VCR_BUCKET_NAME/beta/${'$'}BRANCH_NAME/fixtures/
            fi

            # Cleanup
            rm google-account.json
            gcloud auth application-default revoke
            gcloud auth revoke --all

            echo "Finished"
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}