# White listing IP's using Cloud Armor

This is an example of setting up a project to take advantage of one of the [Cloud Armor features](https://cloud.google.com/armor/) that allows whitelisting of traffic to a compute instance based on ip address. It will set up a single compute instance running nginx that is accessible via a load balanced pool that is managed by cloud armor security policies.

To run the example:
* Set up a Google Cloud Platform account with the [compute engine api enabled](https://console.cloud.google.com/apis/library/compute.googleapis.com)
* [Configure the Google Cloud Provider credentials](https://www.terraform.io/docs/providers/google/index.html#credentials)
* Update the `variables.tf` OR provide overrides in the command line
* Run with a command similar to:
```
terraform apply \
	-var="region=us-west1" \
	-var="region_zone=us-west1-a" \
	-var="project_name=my-project-id-123" \
```

After running `terraform apply` the external ip address of the load balancer will be output to the console. Either enter the ip address into the browser directly or add it to the hosts file on your machine so that it can be accessed at 'mysite.com'.

Navigating to the address either way should result in a 403 rejection. Change the ip address in the whitelist variable in `variables.tf` to your computer's local ip address and re-run `terraform apply` to be able to hit the nginx welcome page on the instance. After the policy has been updated it will need to be propagated to the load balancers which can take up to a few minutes to apply.
