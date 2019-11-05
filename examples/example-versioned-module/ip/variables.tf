variable "name" {
  description = "A name for the ip address resource"
}

variable "labels" {
  type        = map(string)
  description = "A map of key:value labels to apply to the ip address resource"
  default     = {}
}

locals {
  # This ends up being a boolean
  # 1 if there are any entries
  # 0 otherwise
  has_labels = min(1, length(var.labels))
}
