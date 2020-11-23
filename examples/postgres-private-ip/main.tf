# ------------------------------------------------------------------------------
# Provider
# ------------------------------------------------------------------------------

provider "google" {
  version = "=3.30.0"
  project = var.project
  region  = var.region
}

# ------------------------------------------------------------------------------
# Backend
# ------------------------------------------------------------------------------

terraform {
  backend "remote" {
    organization = "myorg"

    workspaces {
      name = "terraform-google-sql"
    }
  }
}

# ------------------------------------------------------------------------------
# Resources
# ------------------------------------------------------------------------------

resource "random_id" "name" {
  byte_length = 4
}

module "postgres" {
  source       = "../.."
  name         = "test-postgres-${random_id.name.hex}"
  db_name      = var.db_name
  project      = var.project
  region       = var.region
  engine       = var.postgres_version
  machine_type = var.machine_type

  # These together will construct the master_user privileges, i.e.
  # 'master_user_name'@'master_user_host' IDENTIFIED BY 'master_user_password'.
  # These should typically be set as the environment variable TF_VAR_master_user_password, etc.
  # so you don't check these into source control."
  master_user_password    = var.master_user_password
  master_user_name        = var.master_user_name
  master_user_host        = "%"
  private_network         = var.network_selflink
  num_read_replicas       = var.num_read_replicas
  read_replica_zones      = var.read_replica_zones
  enable_failover_replica = var.enable_failover_replica
  master_zone             = var.master_zone
  authorized_networks = [
    {
      name  = "all-inbound" # optional
      value = "10.10.0.0/16"
    }
  ]
}