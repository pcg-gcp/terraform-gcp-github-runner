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
| [google_cloud_run_v2_service.control_plane](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_service) | resource |
| [google_cloud_run_v2_service_iam_binding.control_plane](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_service_iam_binding) | resource |
| [google_cloud_scheduler_job.shutdown](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_scheduler_job) | resource |
| [google_project_iam_member.control_plane](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_secret_manager_secret_iam_member.control_plane](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_iam_member) | resource |
| [google_service_account.control_plane](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_service_account.invoker](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_service_account_iam_member.runner_user](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account_iam_member) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allowed_zones"></a> [allowed\_zones](#input\_allowed\_zones) | The zones to deploy the runner instances into.<br>If not set the runner will be deployed into any zone in the region.<br>Also see use\_strict\_zone\_order | `list(string)` | n/a | yes |
| <a name="input_enable_debug"></a> [enable\_debug](#input\_enable\_debug) | Whether to enable debug mode | `bool` | n/a | yes |
| <a name="input_ephemeral"></a> [ephemeral](#input\_ephemeral) | Whether to use ephemeral runners | `bool` | n/a | yes |
| <a name="input_github_app_id"></a> [github\_app\_id](#input\_github\_app\_id) | The GitHub App ID | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | The image to deploy | `string` | n/a | yes |
| <a name="input_image_tag"></a> [image\_tag](#input\_image\_tag) | The tag of the image to deploy | `string` | n/a | yes |
| <a name="input_instance_template_name"></a> [instance\_template\_name](#input\_instance\_template\_name) | The name of the instance template to use for runner instances | `string` | n/a | yes |
| <a name="input_max_instance_count"></a> [max\_instance\_count](#input\_max\_instance\_count) | The maximum number of instances to run | `number` | n/a | yes |
| <a name="input_max_runner_count"></a> [max\_runner\_count](#input\_max\_runner\_count) | The maximum number of runners that should be deployed at the same time | `number` | n/a | yes |
| <a name="input_min_runner_count"></a> [min\_runner\_count](#input\_min\_runner\_count) | The minimum number of runners that should be deployed at all times | `number` | n/a | yes |
| <a name="input_private_key_secret_id"></a> [private\_key\_secret\_id](#input\_private\_key\_secret\_id) | The secret ID of the private key | `string` | n/a | yes |
| <a name="input_private_key_secret_version"></a> [private\_key\_secret\_version](#input\_private\_key\_secret\_version) | The secret version of the private key | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy all resources into | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy all resources into | `string` | n/a | yes |
| <a name="input_runner_labels"></a> [runner\_labels](#input\_runner\_labels) | Github runner labels | `list(string)` | n/a | yes |
| <a name="input_runner_service_account_id"></a> [runner\_service\_account\_id](#input\_runner\_service\_account\_id) | The service account to run the runner as | `string` | n/a | yes |
| <a name="input_shutdown_attempt_timeout"></a> [shutdown\_attempt\_timeout](#input\_shutdown\_attempt\_timeout) | The timeout for the shutdown attempt | `string` | n/a | yes |
| <a name="input_shutdown_schedule"></a> [shutdown\_schedule](#input\_shutdown\_schedule) | The shutdown schedule in cron format | `string` | n/a | yes |
| <a name="input_shutdown_schedule_timezone"></a> [shutdown\_schedule\_timezone](#input\_shutdown\_schedule\_timezone) | The timezone of the shutdown schedule | `string` | n/a | yes |
| <a name="input_use_jit_config"></a> [use\_jit\_config](#input\_use\_jit\_config) | Whether to use JIT config | `bool` | `false` | no |
| <a name="input_use_org_runners"></a> [use\_org\_runners](#input\_use\_org\_runners) | Whether to use github organization runners | `bool` | `false` | no |
| <a name="input_use_strict_zone_order"></a> [use\_strict\_zone\_order](#input\_use\_strict\_zone\_order) | If this is set to true and allowed\_zones is set the runner will always be deployed in the first available zone in the list unless it is unavailable.<br>If allowed\_zones is not set the first zone returned by the API will be used.<br>Otherwise the runner will be deployed in a random zone either from the allowed\_zones list or from the API. | `bool` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_invoker_service_account"></a> [invoker\_service\_account](#output\_invoker\_service\_account) | n/a |
| <a name="output_invoker_service_account_id"></a> [invoker\_service\_account\_id](#output\_invoker\_service\_account\_id) | n/a |
| <a name="output_service_url"></a> [service\_url](#output\_service\_url) | n/a |
<!-- END_TF_DOCS -->