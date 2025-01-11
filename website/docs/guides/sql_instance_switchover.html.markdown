---
page_title: "Performing a SQL Instance Switchover"
description: |-
  A walkthrough for performing a SQL instance switchover through terraform
---

# Performing a SQL Instance Switchover
This page is a brief walkthrough of performing a switchover through terraform. 

## SQL Server

1. Create a **cross-region** primary and cascadable replica. It is recommended to use deletion_protection to prevent accidental deletions.
```
resource "google_sql_database_instance" "original-primary" {
name = "p1"
region = "us-central1"
deletion_protection = true
instance_type = "CLOUD_SQL_INSTANCE"
replica_names = ["p1-r1"] 
    ...
}
resource "google_sql_database_instance" "original-replica" {
name = "p1-r1"
region = "us-east1"
deletion_protection = true
instance_type = "READ_REPLICA_INSTANCE"
master_instance_name = "p1"
replica_configuration {
    cascadable_replica = true
}
...
}
```

2. Invoke switchover on the replica \
a. Change `instance_type` from `READ_REPLICA_INSTANCE` to `CLOUD_SQL_INSTANCE` \
b. Remove `master_instance_name` \
c. Remove `replica_configuration` \
d. Add current primary's name to the replica's `replica_names` list

```diff 
resource "google_sql_database_instance" "original-replica" {
  name = "p1-r1"
  region = "us-east1"
- instance_type = "READ_REPLICA_INSTANCE"
+ instance_type = "CLOUD_SQL_INSTANCE"

- master_instance_name = "p1"
- replica_configuration {
- cascadable_replica = true
- }
+ replica_names = ["p1"]
  ...  
}
```

3. Update the old primary and run `terraform plan` \
a. Change `instance_type` from `CLOUD_SQL_INSTANCE` to `READ_REPLICA_INSTANCE` \
b. Set `master_instance_name` to the new primary (original replica) \
c. Set `replica_configuration` and indicate this is a `cascadable-replica` \
d. Remove old replica from `replica_names` \
    ~> **NOTE**: Do **not** delete the replica_names field, even if it has no replicas remaining. Set replica_names = [ ] to indicate it having no replicas. \
e. Run `terraform plan` and verify that everything is done in-place (or data will be lost)

```diff
resource "google_sql_database_instance" "original-primary" {
  name = "p1"
  region="us-central1"
- instance_type = "CLOUD_SQL_INSTANCE"
+ instance_type = "READ_REPLICA_INSTANCE"
+ master_instance_name = "p1-r1"
+ replica_configuration 
+   cascadable_replica = true
+ }
- replica_names = ["p1-r1"] 
+ replica_names = [] 
  ...
}
```

#### Plan and verify that:
- `terraform plan` outputs **"0 to add, 0 to destroy"**
- `terraform plan` does not say **"must be replaced"** for any resource
- Every resource **"will be updated in-place"**
- Only the 2 instances involved in switchover have planned changes
- (Recommended) Use `deletion_protection` on instances as a safety measure

## MySQL

1. Create a **cross-region, Enterprise Plus edition** primary and replica. The primary should have backup and binary log enabled.

```
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  # Can be any region.
  region              = "us-east1"
  # Any database version that supports Enterprise Plus edition.
  database_version    = "MYSQL_8_0"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
  settings {
    # Any tier that supports Enterprise Plus edition.
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
  
  # You can add more settings.
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  # Can be any region, but must be different from the primary's region.
  region               = "us-west2"
  # Must be same as the primary's database_version.
  database_version     = "MYSQL_8_0"
  instance_type        = "READ_REPLICA_INSTANCE"
  master_instance_name = google_sql_database_instance.original-primary.name
  
  settings {
    # Any tier that supports Enterprise Plus edition.
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
  
  # You can add more settings.
}
```

2. Designate the replica as DR replica of the primary by adding `replication_cluster.failover_dr_replica_name`.
```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "MYSQL_8_0"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
+  replication_cluster {
+    # Note that the format of the name is "project:instance".
+    # If you want to unset DR replica, put empty string in this field.
+    failover_dr_replica_name = "your-project:your-original-replica"
+  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "MYSQL_8_0"
  instance_type        = "READ_REPLICA_INSTANCE"
  master_instance_name = "your-original-primary"
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
```

3. Invoke switchover on the original replica.

* Change `instance_type` from `READ_REPLICA_INSTANCE` to `CLOUD_SQL_INSTANCE`.
* Remove `master_instance_name`.
* Add original primary's name to the original replica's `replica_names` list and `replication_cluster.failover_dr_replica_name`.
* Enable backup and binary log for original replica.

```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "MYSQL_8_0"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
  replication_cluster {
    failover_dr_replica_name = "your-project:your-original-replica"
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "MYSQL_8_0"
-  instance_type        = "READ_REPLICA_INSTANCE"
+  instance_type        = "CLOUD_SQL_INSTANCE"
-  master_instance_name = "your-original-primary"
+  replica_names        = ["your-original-primary"]

+  replication_cluster {
+    failover_dr_replica_name = "your-project:your-original-primary"
+  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
+    backup_configuration {
+      enabled            = true
+      binary_log_enabled = true
+    }    
  }
}
```

4. Update the original primary and run `terraform plan`.
* Change `instance_type` from `CLOUD_SQL_INSTANCE` to `READ_REPLICA_INSTANCE`.
* Set `master_instance_name` to the new primary (original replica).
* (If `replica_names` is present) Remove original replica from `replica_names`.
  * **NOTE**: Do **not** delete the `replica_names` field, even if it has no replicas remaining. Set `replica_names = [ ]` to indicate it having no replicas.
* Remove original replica from `replication_cluster.failover_dr_replica_name` by setting this field to the empty string.
* Disable backup for original primary (because it became a replica).
* Run `terraform plan` and verify that your configuration matches infrastructure. You should see a message like the following:
  * **`No changes. Your infrastructure matches the configuration.`**

```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "MYSQL_8_0"
-  instance_type        = "CLOUD_SQL_INSTANCE"
+  instance_type        = "READ_REPLICA_INSTANCE"
+  master_instance_name = "your-original-replica"
  
  replication_cluster {
-    failover_dr_replica_name = "your-project:your-original-replica"
+    failover_dr_replica_name = ""
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
-      enabled            = true
+      enabled            = false
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "MYSQL_8_0"
  instance_type        = "CLOUD_SQL_INSTANCE"
  replica_names        = ["your-original-primary"]

  replication_cluster {
    failover_dr_replica_name = "your-project:your-original-primary"
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }    
  }
}
```

## PostgreSQL

1. Create a **cross-region, Enterprise Plus edition** primary and replica. The primary should have backup and PITR enabled.

```
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  # Can be any region.
  region              = "us-east1"
  # Any database version that supports Enterprise Plus edition.
  database_version    = "POSTGRES_12"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
  settings {
    # Any tier that supports Enterprise Plus edition.
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }
  }
  
  # You can add more settings.
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  # Can be any region, but must be different from the primary's region.
  region               = "us-west2"
  # Must be same as the primary's database_version.
  database_version     = "POSTGRES_12"
  instance_type        = "READ_REPLICA_INSTANCE"
  master_instance_name = google_sql_database_instance.original-primary.name
  
  settings {
    # Any tier that supports Enterprise Plus edition.
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
  
  # You can add more settings.
}
```

2. Designate the replica as DR replica of the primary by adding `replication_cluster.failover_dr_replica_name`.
```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "POSTGRES_12"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
+  replication_cluster {
+    # Note that the format of the name is "project:instance".
+    # If you want to unset DR replica, put empty string in this field.
+    failover_dr_replica_name = "your-project:your-original-replica"
+  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "POSTGRES_12"
  instance_type        = "READ_REPLICA_INSTANCE"
  master_instance_name = "your-original-primary"
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
```

3. Invoke switchover on the original replica.

* Change `instance_type` from `READ_REPLICA_INSTANCE` to `CLOUD_SQL_INSTANCE`.
* Remove `master_instance_name`.
* Add original primary's name to the original replica's `replica_names` list and `replication_cluster.failover_dr_replica_name`.
* Enable backup and PITR for original replica.

```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "POSTGRES_12"
  instance_type       = "CLOUD_SQL_INSTANCE"
  
  replication_cluster {
    failover_dr_replica_name = "your-project:your-original-replica"
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "POSTGRES_12"
-  instance_type        = "READ_REPLICA_INSTANCE"
+  instance_type        = "CLOUD_SQL_INSTANCE"
-  master_instance_name = "your-original-primary"
+  replica_names        = ["your-original-primary"]

+  replication_cluster {
+    failover_dr_replica_name = "your-project:your-original-primary"
+  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
+    backup_configuration {
+      enabled                        = true
+      point_in_time_recovery_enabled = true
+    }  
  }
}
```

4. Update the original primary and run `terraform plan`.
* Change `instance_type` from `CLOUD_SQL_INSTANCE` to `READ_REPLICA_INSTANCE`.
* Set `master_instance_name` to the new primary (original replica).
* (If `replica_names` is present) Remove original replica from `replica_names`.
  * **NOTE**: Do **not** delete the `replica_names` field, even if it has no replicas remaining. Set `replica_names = [ ]` to indicate it having no replicas.
* Remove original replica from `replication_cluster.failover_dr_replica_name` by setting this field to the empty string.
* Disable backup and PITR for original primary (because it became a replica).
* Run `terraform plan` and verify that your configuration matches infrastructure. You should see a message like the following:
  * **`No changes. Your infrastructure matches the configuration.`**

```diff
resource "google_sql_database_instance" "original-primary" {
  project             = "your-project"
  name                = "your-original-primary"
  region              = "us-east1"
  database_version    = "POSTGRES_12"
-  instance_type        = "CLOUD_SQL_INSTANCE"
+  instance_type        = "READ_REPLICA_INSTANCE"
+  master_instance_name = "your-original-replica"
  
  replication_cluster {
-    failover_dr_replica_name = "your-project:your-original-replica"
+    failover_dr_replica_name = ""
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
-      enabled            = true
+      enabled            = false
-      point_in_time_recovery_enabled = true
+      point_in_time_recovery_enabled = false
    }
  }
}

resource "google_sql_database_instance" "original-replica" {
  project              = "your-project"
  name                 = "your-original-replica"
  region               = "us-west2"
  database_version     = "POSTGRES_12"
  instance_type        = "CLOUD_SQL_INSTANCE"
  replica_names        = ["your-original-primary"]

  replication_cluster {
    failover_dr_replica_name = "your-project:your-original-primary"
  }
  
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }    
  }
}
```
