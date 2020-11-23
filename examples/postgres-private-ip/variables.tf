# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------

variable "project" {
  description = "The project ID to host the database in."
  type        = string
  default     = "my-gcp-project"
}

variable "region" {
  description = "The region to host the database in."
  type        = string
  default     = "asia-northeast1"
}

variable "master_user_name" {
  description = "The username part for the default user credentials, i.e. 'master_user_name'@'master_user_host' IDENTIFIED BY 'master_user_password'. This should typically be set as the environment variable TF_VAR_master_user_name so you don't check it into source control."
  type        = string
  default     = "nirav"
}

variable "master_user_password" {
  description = "The password part for the default user credentials, i.e. 'master_user_name'@'master_user_host' IDENTIFIED BY 'master_user_password'. This should typically be set as the environment variable TF_VAR_master_user_password so you don't check it into source control."
  type        = string
  default     = "admin"
}

# Note, after a name db instance is used, it cannot be reused for up to one week.
variable "name_prefix" {
  description = "The name prefix for the database instance. Will be appended with a random string. Use lowercase letters, numbers, and hyphens. Start with a letter."
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# Generally, these values won't need to be changed.
# ---------------------------------------------------------------------------------------------------------------------

variable "network_selflink" {
  description = "The name or self_link of the Google Compute Engine network to which the cluster is connected. For Shared VPC, set this to the self link of the shared network."
  default     = "projects/sample-project-23fasf35/global/networks/nirav-devops-vpc-admin-west-shared"
}

variable "postgres_version" {
  description = "The engine version of the database, e.g. `POSTGRES_9_6`. See https://cloud.google.com/sql/docs/db-versions for supported versions."
  type        = string
  default     = "POSTGRES_11"
}

variable "machine_type" {
  description = "The machine type to use, see https://cloud.google.com/sql/pricing for more details"
  type        = string
  default     = "db-custom-1-4096"
}

variable "db_name" {
  description = "Name for the db"
  type        = string
  default     = "test_postgres"
}

variable "num_read_replicas" {
  description = "The number of read replicas to create. Cloud SQL will replicate all data from the master to these replicas, which you can use to horizontally scale read traffic."
  type        = number
  default     = 1
}

variable "read_replica_zones" {
  description = "A list of compute zones where read replicas should be created. List size should match 'num_read_replicas'"
  type        = list(string)
  default     = ["asia-northeast1-b"]
}

variable "enable_failover_replica" {
  description = "Set to true to enable failover replica."
  type        = bool
  default     = true
}

variable "master_zone" {
  description = "Preferred zone for the master instance (e.g. 'us-central1-a'). 'region'. If null, Google will auto-assign a zone."
  type        = string
  default     = "asia-northeast1-b"
}