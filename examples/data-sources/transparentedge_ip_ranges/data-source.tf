terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.2.6
      version = ">=0.2.6"
    }
  }
}

provider "transparentedge" {
  # This data source doesn't require authentication
  auth = false
}

data "transparentedge_ip_ranges" "tcdn_ranges" {}

output "tcdn_ranges" {
  value = data.transparentedge_ip_ranges.tcdn_ranges
}

# Ingest directly to AWS SG
resource "aws_security_group" "tcdn" {
  name = "tcdn"

  ingress {
    from_port        = "80"
    to_port          = "80"
    protocol         = "tcp"
    cidr_blocks      = data.transparentedge_ip_ranges.tcdn_ranges.ipv4_cidr_blocks
    ipv6_cidr_blocks = data.transparentedge_ip_ranges.tcdn_ranges.ipv6_cidr_blocks
  }
}
