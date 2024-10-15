// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package transcoder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccTranscoderJobTemplate_transcoderJobTemplateBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTranscoderJobTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTranscoderJobTemplate_transcoderJobTemplateBasicExample(context),
			},
			{
				ResourceName:            "google_transcoder_job_template.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_template_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccTranscoderJobTemplate_transcoderJobTemplateBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_transcoder_job_template" "default" {
  job_template_id = "tf-test-example-job-template%{random_suffix}"
  location = "us-central1"

  config {
    inputs {
      key = "input0"
    }
    edit_list {
      key               = "atom0"
      inputs            = ["input0"]
      start_time_offset = "0s"
    }
    ad_breaks {
      start_time_offset = "3.500s"
    }
    elementary_streams {
      key = "video-stream0"
      video_stream {
        h264 {
          width_pixels      = 640
          height_pixels     = 360
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 550000
          vbv_fullness_bits = 495000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"
        }
      }
    }
    elementary_streams {
      key = "video-stream1"
      video_stream {
        h264 {
          width_pixels      = 1280
          height_pixels     = 720
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 2500000
          vbv_fullness_bits = 2250000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"

        }
      }
    }
    elementary_streams {
      key = "audio-stream0"
      audio_stream {
        codec             = "aac"
        bitrate_bps       = 64000
        channel_count     = 2
        channel_layout    = ["fl", "fr"]
        sample_rate_hertz = 48000
      }
    }
    mux_streams {
      key                = "sd"
      file_name          = "sd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream0", "audio-stream0"]
    }
    mux_streams {
      key                = "hd"
      file_name          = "hd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream1", "audio-stream0"]
    }
  }
  labels = {
    "label" = "key"
  }
}
`, context)
}

func TestAccTranscoderJobTemplate_transcoderJobTemplateOverlaysExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTranscoderJobTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTranscoderJobTemplate_transcoderJobTemplateOverlaysExample(context),
			},
			{
				ResourceName:            "google_transcoder_job_template.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_template_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccTranscoderJobTemplate_transcoderJobTemplateOverlaysExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_transcoder_job_template" "default" {
  job_template_id = "tf-test-example-job-template%{random_suffix}"
  location = "us-central1"
  config {
    inputs {
      key = "input0"
      uri = "gs://example/example.mp4"
    }
    output {
      uri = "gs://example/outputs/"
    }
    edit_list {
      key               = "atom0"
      inputs            = ["input0"]
      start_time_offset = "0s"
    }
    ad_breaks {
      start_time_offset = "3.500s"
    }
    overlays {
      animations {
        animation_fade {
          fade_type         = "FADE_IN"
          start_time_offset = "1.500s"
          end_time_offset   = "3.500s"
          xy {
            x = 1
            y = 0.5
          }
        }
      }
      image {
        uri = "gs://example/overlay.png"
      }
    }
    elementary_streams {
      key = "video-stream0"
      video_stream {
        h264 {
          width_pixels      = 640
          height_pixels     = 360
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 550000
          vbv_fullness_bits = 495000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"

        }
      }
    }
    elementary_streams {
      key = "video-stream1"
      video_stream {
        h264 {
          width_pixels      = 1280
          height_pixels     = 720
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 2500000
          vbv_fullness_bits = 2250000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"
        }
      }
    }
    elementary_streams {
      key = "audio-stream0"
      audio_stream {
        codec             = "aac"
        bitrate_bps       = 64000
        channel_count     = 2
        channel_layout    = ["fl", "fr"]
        sample_rate_hertz = 48000
      }
    }
    mux_streams {
      key                = "sd"
      file_name          = "sd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream0", "audio-stream0"]
    }
    mux_streams {
      key                = "hd"
      file_name          = "hd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream1", "audio-stream0"]
    }
  }
  labels = {
    "label" = "key"
  }
}
`, context)
}

func TestAccTranscoderJobTemplate_transcoderJobTemplateEncryptionsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTranscoderJobTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTranscoderJobTemplate_transcoderJobTemplateEncryptionsExample(context),
			},
			{
				ResourceName:            "google_transcoder_job_template.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_template_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccTranscoderJobTemplate_transcoderJobTemplateEncryptionsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "encryption_key" {
  secret_id = "tf-test-transcoder-encryption-key%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "encryption_key" {
  secret     = google_secret_manager_secret.encryption_key.name
  secret_data = "4A67F2C1B8E93A4F6D3E7890A1BC23DF"
}

resource "google_transcoder_job_template" "default" {
  job_template_id = "tf-test-example-job-template%{random_suffix}"
  location        = "us-central1"

  config {
    elementary_streams {
      key = "es_video"
      video_stream {
        h264 {
          profile      = "main"
          height_pixels = 600
          width_pixels  = 800
          bitrate_bps  = 1000000
          frame_rate   = 60
        }
      }
    }

    elementary_streams {
      key = "es_audio"
      audio_stream {
        codec        = "aac"
        channel_count = 2
        bitrate_bps  = 160000
      }
    }

    encryptions {
      id = "aes-128"
      secret_manager_key_source {
        secret_version = google_secret_manager_secret_version.encryption_key.name
      }
      drm_systems {
        clearkey {}
      }
      aes128 {}
    }

    encryptions {
      id = "cenc"
      secret_manager_key_source {
        secret_version = google_secret_manager_secret_version.encryption_key.name
      }
      drm_systems {
        widevine {}
      }
      mpeg_cenc {
        scheme = "cenc"
      }
    }

    encryptions {
      id = "cbcs"
      secret_manager_key_source {
        secret_version = google_secret_manager_secret_version.encryption_key.name
      }
      drm_systems {
        widevine {}
      }
      mpeg_cenc {
        scheme = "cbcs"
      }
    }

    mux_streams {
      key                 = "ts_aes128"
      container           = "ts"
      elementary_streams  = ["es_video", "es_audio"]
      segment_settings {
        segment_duration = "6s"
      }
      encryption_id = "aes-128"
    }

    mux_streams {
      key                 = "fmp4_cenc_video"
      container           = "fmp4"
      elementary_streams  = ["es_video"]
      segment_settings {
        segment_duration = "6s"
      }
      encryption_id = "cenc"
    }

    mux_streams {
      key                 = "fmp4_cenc_audio"
      container           = "fmp4"
      elementary_streams  = ["es_audio"]
      segment_settings {
        segment_duration = "6s"
      }
      encryption_id = "cenc"
    }

    mux_streams {
      key                 = "fmp4_cbcs_video"
      container           = "fmp4"
      elementary_streams  = ["es_video"]
      segment_settings {
        segment_duration = "6s"
      }
      encryption_id = "cbcs"
    }

    mux_streams {
      key                 = "fmp4_cbcs_audio"
      container           = "fmp4"
      elementary_streams  = ["es_audio"]
      segment_settings {
        segment_duration = "6s"
      }
      encryption_id = "cbcs"
    }

    manifests {
      file_name = "manifest_aes128.m3u8"
      type      = "HLS"
      mux_streams = ["ts_aes128"]
    }

    manifests {
      file_name = "manifest_cenc.mpd"
      type      = "DASH"
      mux_streams = ["fmp4_cenc_video", "fmp4_cenc_audio"]
    }

    manifests {
      file_name = "manifest_cbcs.mpd"
      type      = "DASH"
      mux_streams = ["fmp4_cbcs_video", "fmp4_cbcs_audio"]
    }
  }
  labels = {
    "label" = "key"
  }
}
`, context)
}

func TestAccTranscoderJobTemplate_transcoderJobTemplatePubsubExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTranscoderJobTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTranscoderJobTemplate_transcoderJobTemplatePubsubExample(context),
			},
			{
				ResourceName:            "google_transcoder_job_template.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"job_template_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccTranscoderJobTemplate_transcoderJobTemplatePubsubExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "transcoder_notifications" {
  name = "tf-test-transcoder-notifications%{random_suffix}"
}

resource "google_transcoder_job_template" "default" {
  job_template_id = "tf-test-example-job-template%{random_suffix}"
  location = "us-central1"
  config {
    inputs {
      key = "input0"
      uri = "gs://example/example.mp4"
    }
    output {
      uri = "gs://example/outputs/"
    }
    edit_list {
      key               = "atom0"
      inputs            = ["input0"]
      start_time_offset = "0s"
    }
    ad_breaks {
      start_time_offset = "3.500s"
    }
    elementary_streams {
      key = "video-stream0"
      video_stream {
        h264 {
          width_pixels      = 640
          height_pixels     = 360
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 550000
          vbv_fullness_bits = 495000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"

        }
      }
    }
    elementary_streams {
      key = "video-stream1"
      video_stream {
        h264 {
          width_pixels      = 1280
          height_pixels     = 720
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 2500000
          vbv_fullness_bits = 2250000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"
        }
      }
    }
    elementary_streams {
      key = "audio-stream0"
      audio_stream {
        codec             = "aac"
        bitrate_bps       = 64000
        channel_count     = 2
        channel_layout    = ["fl", "fr"]
        sample_rate_hertz = 48000
      }
    }
    mux_streams {
      key                = "sd"
      file_name          = "sd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream0", "audio-stream0"]
    }
    mux_streams {
      key                = "hd"
      file_name          = "hd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream1", "audio-stream0"]
    }
    pubsub_destination {
      topic = google_pubsub_topic.transcoder_notifications.id
    }
  }
  labels = {
    "label" = "key"
  }
}
`, context)
}

func testAccCheckTranscoderJobTemplateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_transcoder_job_template" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{TranscoderBasePath}}projects/{{project}}/locations/{{location}}/jobTemplates/{{job_template_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TranscoderJobTemplate still exists at %s", url)
			}
		}

		return nil
	}
}
