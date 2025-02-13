<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.8 |
| <a name="requirement_google"></a> [google](#requirement\_google) | ~> 6.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | ~> 6.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_cloud_run_v2_service.webhook](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_service) | resource |
| [google_service_account.webhook](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/service_account) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_control_plane_url"></a> [control\_plane\_url](#input\_control\_plane\_url) | The control plane URL | `string` | n/a | yes |
| <a name="input_enable_debug"></a> [enable\_debug](#input\_enable\_debug) | Whether to enable debug mode | `bool` | n/a | yes |
| <a name="input_forward_delay_seconds"></a> [forward\_delay\_seconds](#input\_forward\_delay\_seconds) | The number of seconds the webhook handler delays events before forwarding them to the control plane | `number` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | The image to deploy | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | The tag of the image to deploy | `string` | n/a | yes |
| <a name="input_invoker_service_account"></a> [invoker\_service\_account](#input\_invoker\_service\_account) | The service account to trigger the control plane | `string` | n/a | yes |
| <a name="input_invoker_service_account_id"></a> [invoker\_service\_account\_id](#input\_invoker\_service\_account\_id) | The service account ID to trigger the control plane | `string` | n/a | yes |
| <a name="input_max_instance_count"></a> [max\_instance\_count](#input\_max\_instance\_count) | The maximum number of instances to run | `number` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy all resources into | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy all resources into | `string` | `"europe-west1"` | no |
| <a name="input_runner_labels"></a> [runner\_labels](#input\_runner\_labels) | Github runner labels | `list(string)` | n/a | yes |
| <a name="input_task_queue_path"></a> [task\_queue\_path](#input\_task\_queue\_path) | The path to the task queue | `string` | n/a | yes |
| <a name="input_webhook_secret_id"></a> [webhook\_secret\_id](#input\_webhook\_secret\_id) | The secret ID for the webhook | `string` | n/a | yes |
| <a name="input_webhook_secret_version"></a> [webhook\_secret\_version](#input\_webhook\_secret\_version) | The secret version for the webhook | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_webhook_url"></a> [webhook\_url](#output\_webhook\_url) | n/a |
<!-- END_TF_DOCS -->