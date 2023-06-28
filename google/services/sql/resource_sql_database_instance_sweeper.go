// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

func init() {
	sweeper.AddTestSweepers("SQLDatabaseInstance", testSweepSQLDatabaseInstance)
}

func testSweepSQLDatabaseInstance(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	found, err := config.NewSqlAdminClient(config.UserAgent).Instances.List(config.Project).Do()
	if err != nil {
		log.Printf("error listing databases: %s", err)
		return nil
	}

	if len(found.Items) == 0 {
		log.Printf("No databases found")
		return nil
	}

	running := map[string]struct{}{}

	for _, d := range found.Items {
		if !sweeper.IsSweepableTestResource(d.Name) {
			continue
		}

		if d.State != "RUNNABLE" {
			continue
		}
		running[d.Name] = struct{}{}
	}

	for _, d := range found.Items {
		if !sweeper.IsSweepableTestResource(d.Name) {
			continue
		}

		// don't delete replicas, we'll take care of that
		// when deleting the database they replicate
		if d.ReplicaConfiguration != nil {
			continue
		}
		log.Printf("Destroying SQL Instance (%s)", d.Name)

		// replicas need to be stopped and destroyed before destroying a master
		// instance. The ordering slice tracks replica databases for a given master
		// and we call destroy on them before destroying the master
		var ordering []string
		for _, replicaName := range d.ReplicaNames {
			// don't try to stop replicas that aren't running
			if _, ok := running[replicaName]; !ok {
				ordering = append(ordering, replicaName)
				continue
			}

			// need to stop replication before being able to destroy a database
			op, err := config.NewSqlAdminClient(config.UserAgent).Instances.StopReplica(config.Project, replicaName).Do()

			if err != nil {
				log.Printf("error, failed to stop replica instance (%s) for instance (%s): %s", replicaName, d.Name, err)
				return nil
			}

			err = SqlAdminOperationWaitTime(config, op, config.Project, "Stop Replica", config.UserAgent, 10*time.Minute)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("Replication operation not found")
				} else {
					log.Printf("Error waiting for sqlAdmin operation: %s", err)
					return nil
				}
			}

			ordering = append(ordering, replicaName)
		}

		// ordering has a list of replicas (or none), now add the primary to the end
		ordering = append(ordering, d.Name)

		for _, db := range ordering {
			// destroy instances, replicas first
			op, err := config.NewSqlAdminClient(config.UserAgent).Instances.Delete(config.Project, db).Do()

			if err != nil {
				if strings.Contains(err.Error(), "409") {
					// the GCP api can return a 409 error after the delete operation
					// reaches a successful end
					log.Printf("Operation not found, got 409 response")
					continue
				}

				log.Printf("Error, failed to delete instance %s: %s", db, err)
				return nil
			}

			err = SqlAdminOperationWaitTime(config, op, config.Project, "Delete Instance", config.UserAgent, 10*time.Minute)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("SQL instance not found")
					continue
				}
				log.Printf("Error, failed to delete instance %s: %s", db, err)
				return nil
			}
		}
	}

	return nil
}
