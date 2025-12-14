# provider.tf

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0" 
    }
  }
}

# Configure the Google Cloud provider using variables
provider "google" {
  project = var.project_id
  region  = var.region 
}