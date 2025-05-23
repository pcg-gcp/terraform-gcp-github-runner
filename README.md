# terraform-gcp-github-runner

This repository provides a Terraform module to deploy self-hosted GitHub Runners
on Google Cloud GCE Instances. These runners scale automatically based on GitHub
events and even support scale down to zero and spot instances.

## License

This module is licensed under the Apache License Version 2.0.
Please see the LICENSE file for more details.

## Contributing

This repository uses [Conventional Commit messages](https://www.conventionalcommits.org/)
to automatically generate a changelog and versioning. The following commit types
are supported:

- fix: A bug fix (PATCH)
- feat: A new feature (MINOR)
- fix!, feat!, refactor! etc.: A breaking change (MAJOR)

Please make sure to follow the commit message convention when creating pull requests.

<details>
<summary>Terraform documentation</summary>
<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.9 |
| <a name="requirement_google"></a> [google](#requirement\_google) | ~> 6.0 |
| <a name="requirement_http"></a> [http](#requirement\_http) | ~> 3.0 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | ~> 6.0 |
| <a name="provider_random"></a> [random](#provider\_random) | ~> 3.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_artifact_registry"></a> [artifact\_registry](#module\_artifact\_registry) | ./modules/artifact_registry | n/a |
| <a name="module_control_plane"></a> [control\_plane](#module\_control\_plane) | ./modules/control_plane | n/a |
| <a name="module_runner_template"></a> [runner\_template](#module\_runner\_template) | ./modules/runner_template | n/a |
| <a name="module_webhook"></a> [webhook](#module\_webhook) | ./modules/webhook | n/a |

## Resources

| Name | Type |
|------|------|
| [google_cloud_tasks_queue.github_events](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_tasks_queue) | resource |
| [google_project_service.required_services](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_service) | resource |
| [google_secret_manager_secret.github_auth_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret) | resource |
| [google_secret_manager_secret.webhook_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret) | resource |
| [google_secret_manager_secret_version.github_auth_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_version) | resource |
| [google_secret_manager_secret_version.webhook_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_version) | resource |
| [random_string.queue_suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_allowed_zones"></a> [allowed\_zones](#input\_allowed\_zones) | The zones to deploy the runner instances into.<br/>If not set the runner will be deployed into any zone in the region.<br/>Also see use\_strict\_zone\_order | `list(string)` | `[]` | no |
| <a name="input_control_plane_account_id"></a> [control\_plane\_account\_id](#input\_control\_plane\_account\_id) | Account id of the control plane service account only used if service account management is disabled. | `string` | `"ghr-control-plane"` | no |
| <a name="input_control_plane_oci_image"></a> [control\_plane\_oci\_image](#input\_control\_plane\_oci\_image) | The control plane OCI image to deploy. This needs to be the full path without the image tag. This is only used if use\_remote\_repository is set to false | `string` | `""` | no |
| <a name="input_control_plane_version"></a> [control\_plane\_version](#input\_control\_plane\_version) | The version of the control plane to deploy. If not set the module version will be used | `string` | `""` | no |
| <a name="input_disable_service_account_management"></a> [disable\_service\_account\_management](#input\_disable\_service\_account\_management) | Whether the used service accounts should be create by this module. | `bool` | `false` | no |
| <a name="input_enable_apis"></a> [enable\_apis](#input\_enable\_apis) | Whether this module should enable the required apis. | `bool` | `true` | no |
| <a name="input_enable_debug"></a> [enable\_debug](#input\_enable\_debug) | Whether to enable debug mode | `bool` | `false` | no |
| <a name="input_ephemeral"></a> [ephemeral](#input\_ephemeral) | Whether to use ephemeral runners | `bool` | `false` | no |
| <a name="input_forward_delay_seconds"></a> [forward\_delay\_seconds](#input\_forward\_delay\_seconds) | The number of seconds the webhook handler delays events before forwarding them to the control plane | `number` | `10` | no |
| <a name="input_github_app_id"></a> [github\_app\_id](#input\_github\_app\_id) | The GitHub App ID | `string` | n/a | yes |
| <a name="input_github_app_private_key_base64"></a> [github\_app\_private\_key\_base64](#input\_github\_app\_private\_key\_base64) | The base64 encoded private key of the GitHub App | `string` | n/a | yes |
| <a name="input_grant_runner_user_sudo"></a> [grant\_runner\_user\_sudo](#input\_grant\_runner\_user\_sudo) | Whether to grant the runner user sudo access. This only has an effect if include\_install\_step is set to true | `bool` | `false` | no |
| <a name="input_include_install_step"></a> [include\_install\_step](#input\_include\_install\_step) | Whether to include the install step for the setup script | `bool` | `true` | no |
| <a name="input_invoker_account_id"></a> [invoker\_account\_id](#input\_invoker\_account\_id) | Account id of the invoker service account only used if service account management is disabled. | `string` | `"ghr-cp-invoker"` | no |
| <a name="input_max_runner_count"></a> [max\_runner\_count](#input\_max\_runner\_count) | The maximum number of runners that should be deployed at the same time | `number` | `10` | no |
| <a name="input_min_runner_count"></a> [min\_runner\_count](#input\_min\_runner\_count) | The minimum number of runners that should be deployed at all times | `number` | `0` | no |
| <a name="input_node_version"></a> [node\_version](#input\_node\_version) | NodeJS version to install | `string` | `""` | no |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy all resources into | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy all resources into | `string` | `"europe-west1"` | no |
| <a name="input_remote_control_plane_image_name"></a> [remote\_control\_plane\_image\_name](#input\_remote\_control\_plane\_image\_name) | The name of the control plane image in the remote repository. This should only be the image name as it is combined with the repository url to form the full path. This is only used if use\_remote\_repository is set to true | `string` | `"pcg-gcp/terraform-gcp-github-runner/control-plane"` | no |
| <a name="input_remote_repository_url"></a> [remote\_repository\_url](#input\_remote\_repository\_url) | The URL of the remote repository to clone. This is only used if use\_remote\_repository is set to true | `string` | `"https://ghcr.io"` | no |
| <a name="input_remote_webhook_image_name"></a> [remote\_webhook\_image\_name](#input\_remote\_webhook\_image\_name) | The name of the webhook image in the remote repository. This should only be the image name as it is combined with the repository url to form the full path. This is only used if use\_remote\_repository is set to true | `string` | `"pcg-gcp/terraform-gcp-github-runner/webhook"` | no |
| <a name="input_runner_account_id"></a> [runner\_account\_id](#input\_runner\_account\_id) | Account id of the runner service account only used if service account management is disabled. | `string` | `"ghr-runner"` | no |
| <a name="input_runner_additional_disks"></a> [runner\_additional\_disks](#input\_runner\_additional\_disks) | Additional disks to attach to the runner | `list(any)` | `[]` | no |
| <a name="input_runner_dir"></a> [runner\_dir](#input\_runner\_dir) | The directory to run the runner in | `string` | `"/opt/github-runner"` | no |
| <a name="input_runner_disk_size_gb"></a> [runner\_disk\_size\_gb](#input\_runner\_disk\_size\_gb) | The disk size in GB to deploy | `number` | `null` | no |
| <a name="input_runner_disk_type"></a> [runner\_disk\_type](#input\_runner\_disk\_type) | The disk type to deploy | `string` | `"pd-balanced"` | no |
| <a name="input_runner_extra_labels"></a> [runner\_extra\_labels](#input\_runner\_extra\_labels) | Github runner extra labels. These should not include github read-only labels like 'self-hosted' or 'linux' | `list(string)` | `[]` | no |
| <a name="input_runner_image_path"></a> [runner\_image\_path](#input\_runner\_image\_path) | The image to deploy | `string` | `"ubuntu-os-cloud/ubuntu-2404-lts-amd64"` | no |
| <a name="input_runner_machine_type"></a> [runner\_machine\_type](#input\_runner\_machine\_type) | The machine type to deploy | `string` | n/a | yes |
| <a name="input_runner_on_host_maintenance"></a> [runner\_on\_host\_maintenance](#input\_runner\_on\_host\_maintenance) | The maintenance policy for the runner | `string` | `"MIGRATE"` | no |
| <a name="input_runner_use_spot_vms"></a> [runner\_use\_spot\_vms](#input\_runner\_use\_spot\_vms) | Whether to use spot VMs for the runner | `bool` | `false` | no |
| <a name="input_runner_user"></a> [runner\_user](#input\_runner\_user) | The user to run the runner as | `string` | `"ghrunner"` | no |
| <a name="input_runner_version"></a> [runner\_version](#input\_runner\_version) | GitHub Runner version to install | `string` | `""` | no |
| <a name="input_shutdown_attempt_timeout"></a> [shutdown\_attempt\_timeout](#input\_shutdown\_attempt\_timeout) | The timeout for the shutdown attempt | `string` | `"320s"` | no |
| <a name="input_shutdown_schedule"></a> [shutdown\_schedule](#input\_shutdown\_schedule) | The shutdown schedule in cron format | `string` | `"*/5 * * * *"` | no |
| <a name="input_shutdown_schedule_timezone"></a> [shutdown\_schedule\_timezone](#input\_shutdown\_schedule\_timezone) | The timezone of the shutdown schedule | `string` | `"Etc/UTC"` | no |
| <a name="input_subnet_name"></a> [subnet\_name](#input\_subnet\_name) | The subnet to deploy runner instances into | `string` | n/a | yes |
| <a name="input_use_jit_config"></a> [use\_jit\_config](#input\_use\_jit\_config) | Whether to use JIT config | `bool` | `false` | no |
| <a name="input_use_org_runners"></a> [use\_org\_runners](#input\_use\_org\_runners) | Whether to use github organization runners | `bool` | `false` | no |
| <a name="input_use_remote_repository"></a> [use\_remote\_repository](#input\_use\_remote\_repository) | Whether to use a remote repository for caching images | `bool` | `true` | no |
| <a name="input_use_strict_zone_order"></a> [use\_strict\_zone\_order](#input\_use\_strict\_zone\_order) | If this is set to true and allowed\_zones is set the runner will always be deployed in the first available zone in the list unless it is unavailable.<br/>If allowed\_zones is not set the first zone returned by the API will be used.<br/>Otherwise the runner will be deployed in a random zone either from the allowed\_zones list or from the API. | `bool` | `false` | no |
| <a name="input_vpc_name"></a> [vpc\_name](#input\_vpc\_name) | The VPC to deploy runner instances into | `string` | n/a | yes |
| <a name="input_webhook_account_id"></a> [webhook\_account\_id](#input\_webhook\_account\_id) | Account id of the webhook service account only used if service account management is disabled. | `string` | `"ghr-webhook"` | no |
| <a name="input_webhook_oci_image"></a> [webhook\_oci\_image](#input\_webhook\_oci\_image) | The webhook OCI image to deploy. This needs to be the full path withouth the image tag. This is only used if use\_remote\_repository is set to false | `string` | `""` | no |
| <a name="input_webhook_secret"></a> [webhook\_secret](#input\_webhook\_secret) | The secret to use | `string` | n/a | yes |
| <a name="input_webhook_version"></a> [webhook\_version](#input\_webhook\_version) | The version of the webhook to deploy. If not set the module version will be used | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_webhook_secret"></a> [webhook\_secret](#output\_webhook\_secret) | n/a |
| <a name="output_webhook_url"></a> [webhook\_url](#output\_webhook\_url) | n/a |
<!-- END_TF_DOCS -->
</details>
