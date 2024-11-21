---
page_title: "Performing a SQL Instance Switchover"
description: |-
  A walkthrough for performing a SQL instance switchover through terraform
---

# Performing a SQL Instance Switchover
This page is a brief walkthrough of performing a switchover through terraform. 

  ~> **NOTE:** Only supported for SQL Server.

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