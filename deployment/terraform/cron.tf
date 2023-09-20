resource "kubernetes_cron_job_v1" "cron_fill" {
  metadata {
    name = var.gke_cron_fill_name
    labels = {
      app = var.gke_cron_fill_name
    }
  }

  spec {
    schedule = var.gke_cron_fill_schedule
    job_template {
      metadata {
        labels = {
          app = var.gke_cron_fill_name
        }
      }
      spec {
        template {
          metadata {
            labels = {
              app = var.gke_cron_fill_name
            }
          }
          spec {
            container {
              name    = var.gke_cron_fill_name
              image   = var.gcr_image_name
              command = ["./shimakaze"]
              args    = ["cron", "fill"]
              env {
                name  = "SHIMAKAZE_CACHE_DIALECT"
                value = var.shimakaze_cache_dialect
              }
              env {
                name  = "SHIMAKAZE_CACHE_ADDRESS"
                value = var.shimakaze_cache_address
              }
              env {
                name  = "SHIMAKAZE_CACHE_PASSWORD"
                value = var.shimakaze_cache_password
              }
              env {
                name  = "SHIMAKAZE_CACHE_TIME"
                value = var.shimakaze_cache_time
              }
              env {
                name  = "SHIMAKAZE_DB_ADDRESS"
                value = var.shimakaze_db_address
              }
              env {
                name  = "SHIMAKAZE_DB_NAME"
                value = var.shimakaze_db_name
              }
              env {
                name  = "SHIMAKAZE_DB_USER"
                value = var.shimakaze_db_user
              }
              env {
                name  = "SHIMAKAZE_DB_PASSWORD"
                value = var.shimakaze_db_password
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_DIALECT"
                value = var.shimakaze_pubsub_dialect
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_ADDRESS"
                value = var.shimakaze_pubsub_address
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_PASSWORD"
                value = var.shimakaze_pubsub_password
              }
              env {
                name  = "SHIMAKAZE_CRON_UPDATE_LIMIT"
                value = var.shimakaze_cron_update_limit
              }
              env {
                name  = "SHIMAKAZE_CRON_FILL_LIMIT"
                value = var.shimakaze_cron_fill_limit
              }
              env {
                name  = "SHIMAKAZE_CRON_AGENCY_AGE"
                value = var.shimakaze_cron_agency_age
              }
              env {
                name  = "SHIMAKAZE_CRON_ACTIVE_AGE"
                value = var.shimakaze_cron_active_age
              }
              env {
                name  = "SHIMAKAZE_CRON_RETIRED_AGE"
                value = var.shimakaze_cron_retired_age
              }
              env {
                name  = "SHIMAKAZE_YOUTUBE_KEY"
                value = var.shimakaze_youtube_key
              }
              env {
                name  = "SHIMAKAZE_TWITCH_CLIENT_ID"
                value = var.shimakaze_twitch_client_id
              }
              env {
                name  = "SHIMAKAZE_TWITCH_CLIENT_SECRET"
                value = var.shimakaze_twitch_client_secret
              }
              env {
                name  = "SHIMAKAZE_LOG_JSON"
                value = var.shimakaze_log_json
              }
              env {
                name  = "SHIMAKAZE_LOG_LEVEL"
                value = var.shimakaze_log_level
              }
              env {
                name  = "SHIMAKAZE_NEWRELIC_LICENSE_KEY"
                value = var.shimakaze_newrelic_license_key
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_cron_job_v1" "cron_update" {
  metadata {
    name = var.gke_cron_update_name
    labels = {
      app = var.gke_cron_update_name
    }
  }

  spec {
    schedule = var.gke_cron_update_schedule
    job_template {
      metadata {
        labels = {
          app = var.gke_cron_update_name
        }
      }
      spec {
        template {
          metadata {
            labels = {
              app = var.gke_cron_update_name
            }
          }
          spec {
            container {
              name    = var.gke_cron_update_name
              image   = var.gcr_image_name
              command = ["./shimakaze"]
              args    = ["cron", "update"]
              env {
                name  = "SHIMAKAZE_CACHE_DIALECT"
                value = var.shimakaze_cache_dialect
              }
              env {
                name  = "SHIMAKAZE_CACHE_ADDRESS"
                value = var.shimakaze_cache_address
              }
              env {
                name  = "SHIMAKAZE_CACHE_PASSWORD"
                value = var.shimakaze_cache_password
              }
              env {
                name  = "SHIMAKAZE_CACHE_TIME"
                value = var.shimakaze_cache_time
              }
              env {
                name  = "SHIMAKAZE_DB_ADDRESS"
                value = var.shimakaze_db_address
              }
              env {
                name  = "SHIMAKAZE_DB_NAME"
                value = var.shimakaze_db_name
              }
              env {
                name  = "SHIMAKAZE_DB_USER"
                value = var.shimakaze_db_user
              }
              env {
                name  = "SHIMAKAZE_DB_PASSWORD"
                value = var.shimakaze_db_password
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_DIALECT"
                value = var.shimakaze_pubsub_dialect
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_ADDRESS"
                value = var.shimakaze_pubsub_address
              }
              env {
                name  = "SHIMAKAZE_PUBSUB_PASSWORD"
                value = var.shimakaze_pubsub_password
              }
              env {
                name  = "SHIMAKAZE_CRON_UPDATE_LIMIT"
                value = var.shimakaze_cron_update_limit
              }
              env {
                name  = "SHIMAKAZE_CRON_FILL_LIMIT"
                value = var.shimakaze_cron_fill_limit
              }
              env {
                name  = "SHIMAKAZE_CRON_AGENCY_AGE"
                value = var.shimakaze_cron_agency_age
              }
              env {
                name  = "SHIMAKAZE_CRON_ACTIVE_AGE"
                value = var.shimakaze_cron_active_age
              }
              env {
                name  = "SHIMAKAZE_CRON_RETIRED_AGE"
                value = var.shimakaze_cron_retired_age
              }
              env {
                name  = "SHIMAKAZE_YOUTUBE_KEY"
                value = var.shimakaze_youtube_key
              }
              env {
                name  = "SHIMAKAZE_TWITCH_CLIENT_ID"
                value = var.shimakaze_twitch_client_id
              }
              env {
                name  = "SHIMAKAZE_TWITCH_CLIENT_SECRET"
                value = var.shimakaze_twitch_client_secret
              }
              env {
                name  = "SHIMAKAZE_LOG_JSON"
                value = var.shimakaze_log_json
              }
              env {
                name  = "SHIMAKAZE_LOG_LEVEL"
                value = var.shimakaze_log_level
              }
              env {
                name  = "SHIMAKAZE_NEWRELIC_LICENSE_KEY"
                value = var.shimakaze_newrelic_license_key
              }
            }
          }
        }
      }
    }
  }
}
