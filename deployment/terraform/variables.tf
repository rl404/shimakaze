variable "gcp_project_id" {
  type        = string
  description = "GCP project id"
}

variable "gcp_region" {
  type        = string
  description = "GCP project region"
}

variable "gke_cluster_name" {
  type        = string
  description = "GKE cluster name"
}

variable "gke_location" {
  type        = string
  description = "GKE location"
}

variable "gke_pool_name" {
  type        = string
  description = "GKE node pool name"
}

variable "gke_node_preemptible" {
  type        = bool
  description = "GKE node preemptible"
}

variable "gke_node_machine_type" {
  type        = string
  description = "GKE node machine type"
}

variable "gcr_image_name" {
  type        = string
  description = "GCR image name"
}

variable "gke_deployment_consumer_name" {
  type        = string
  description = "GKE deployment consumer name"
}

variable "gke_cron_fill_name" {
  type        = string
  description = "GKE cron fill name"
}

variable "gke_cron_fill_schedule" {
  type        = string
  description = "GKE cron fill schedule"
}

variable "gke_cron_update_name" {
  type        = string
  description = "GKE cron update name"
}

variable "gke_cron_update_schedule" {
  type        = string
  description = "GKE cron update schedule"
}

variable "cloud_run_name" {
  type        = string
  description = "Google cloud run name"
}

variable "cloud_run_location" {
  type        = string
  description = "Google cloud run location"
}

variable "shimakaze_cache_dialect" {
  type        = string
  description = "Cache dialect"
}

variable "shimakaze_cache_address" {
  type        = string
  description = "Cache address"
}

variable "shimakaze_cache_password" {
  type        = string
  description = "Cache password"
}

variable "shimakaze_cache_time" {
  type        = string
  description = "Cache time"
}

variable "shimakaze_db_address" {
  type        = string
  description = "Database address"
}

variable "shimakaze_db_name" {
  type        = string
  description = "Database name"
}

variable "shimakaze_db_user" {
  type        = string
  description = "Database user"
}

variable "shimakaze_db_password" {
  type        = string
  description = "Database password"
}

variable "shimakaze_pubsub_dialect" {
  type        = string
  description = "Pubsub dialect"
}

variable "shimakaze_pubsub_address" {
  type        = string
  description = "Pubsub address"
}

variable "shimakaze_pubsub_password" {
  type        = string
  description = "Pubsub password"
}

variable "shimakaze_cron_update_limit" {
  type        = number
  description = "Cron update limit"
}

variable "shimakaze_cron_fill_limit" {
  type        = number
  description = "Cron fill limit"
}

variable "shimakaze_cron_agency_age" {
  type        = number
  description = "Cron agency age"
}

variable "shimakaze_cron_active_age" {
  type        = number
  description = "Cron active age"
}

variable "shimakaze_cron_retired_age" {
  type        = number
  description = "Cron retired age"
}

variable "shimakaze_youtube_key" {
  type        = string
  description = "Youtube API key"
}

variable "shimakaze_twitch_client_id" {
  type        = string
  description = "Twitch client id"
}

variable "shimakaze_twitch_client_secret" {
  type        = string
  description = "Twitch client secret"
}

variable "shimakaze_log_json" {
  type        = bool
  description = "Log json"
}

variable "shimakaze_log_level" {
  type        = number
  description = "Log level"
}

variable "shimakaze_newrelic_license_key" {
  type        = string
  description = "Newrelic license key"
}
