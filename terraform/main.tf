# main.tf 

data "google_compute_network" "vpc_network" {
  name = var.network_name
}

# Define the GCP VPC Firewall Rule resource
resource "google_compute_firewall" "allow_custom_ports_ingress" { 
  name    = "allow-custom-ports-ingress-to-${var.target_tag}"
  network = data.google_compute_network.vpc_network.self_link  

  direction = "INGRESS"
  priority = 100
  target_tags = [var.target_tag]
  source_ranges = [var.allowed_source_cidr]

  # The protocol and port that are permitted
  allow {
    protocol = "tcp"
    # UPDATED: Allow TCP ports 9000 and 9001
    ports    = ["9000", "9001"] 
  }
   
  description = "Allow custom ports (TCP:9000, 9001) ingress from specified CIDR to target tagged instances"
}

output "firewall_rule_name" {
  value = google_compute_firewall.allow_custom_ports_ingress.name
}