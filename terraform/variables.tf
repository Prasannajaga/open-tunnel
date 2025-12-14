# variables.tf

variable "project_id" {
  description = "The ID of the GCP project where the firewall will be created"
  type        = string
}

variable "region" {
  description = "The region for the provider configuration (Firewall rules are global, but this is good practice)"
  type        = string 
}

variable "network_name" {
  description = "The name of the VPC network to apply the firewall rule to (e.g., 'default' or a custom VPC name)"
  type        = string
}

variable "target_tag" {
  description = "The network tag that identifies the target instances (e.g., 'web-server')"
  type        = string
}

variable "allowed_source_cidr" {
  description = "The source IP CIDR range that is allowed to connect (e.g., '0.0.0.0/0' for all, or '203.0.113.0/24')"
  type        = string
}