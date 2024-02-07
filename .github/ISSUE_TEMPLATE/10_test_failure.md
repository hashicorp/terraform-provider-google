---
name: Test failure
labels: test failure
about: (Internal) For reporting test failures on the nightly builds
title: "Failing test(s): TestAccWhatever"

---
<!--- This is a template for reporting test failures on nightly builds. It should only be used by core contributors who have access to our CI/CD results. --->

### Failure rates

- 100% since YYYY-MM-DD
- X% in MONTH

<!--
TeamCity projects can be found here: https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud?mode=builds
-->



### Impacted tests
<!-- List all impacted tests for searchability. The title of the issue can instead list one or more groups of tests, or describe the overall root cause. -->

- TestAccWhatever

### Affected Resource(s)

<!--- List the affected resources and data sources. Use google_* if all resources or data sources are affected. --->

* google_XXXXX

### Nightly build test history

<!-- Link to the test failure(s) page ie https://hashicorp.teamcity.com/test/1458901879521596451?currentProjectId=TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS&expandTestHistoryChartSection=true -->
- Link

<!-- The error message that displays in the tests tab, for reference -->
### Message(s)

```

```
