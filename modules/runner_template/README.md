<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.8 |
| <a name="requirement_google"></a> [google](#requirement\_google) | ~> 5.0 |
| <a name="requirement_http"></a> [http](#requirement\_http) | ~> 3.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | ~> 5.0 |
| <a name="provider_http"></a> [http](#provider\_http) | ~> 3.0 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_compute_instance_template.runner](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance_template) | resource |
| [google_project_iam_member.runner](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_iam_member) | resource |
| [google_service_account.runner](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/service_account) | resource |
| [google_storage_bucket.runner_bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket) | resource |
| [google_storage_bucket_iam_member.runner](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_iam_member) | resource |
| [google_storage_bucket_object.startup_script](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_object) | resource |
| [random_string.bucket_suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [http_http.github_runner_release_json](https://registry.terraform.io/providers/hashicorp/http/latest/docs/data-sources/http) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additional_disks"></a> [additional\_disks](#input\_additional\_disks) | Additional disks to attach to the runner | `list(any)` | n/a | yes |
| <a name="input_disk_size_gb"></a> [disk\_size\_gb](#input\_disk\_size\_gb) | The disk size in GB to deploy | `number` | n/a | yes |
| <a name="input_disk_type"></a> [disk\_type](#input\_disk\_type) | The disk type to deploy | `string` | n/a | yes |
| <a name="input_image_path"></a> [image\_path](#input\_image\_path) | The image to deploy | `string` | n/a | yes |
| <a name="input_include_install_step"></a> [include\_install\_step](#input\_include\_install\_step) | Whether to include the install step for the setup script | `bool` | n/a | yes |
| <a name="input_machine_type"></a> [machine\_type](#input\_machine\_type) | The machine type to deploy | `string` | n/a | yes |
| <a name="input_node_version"></a> [node\_version](#input\_node\_version) | NodeJS version to install | `string` | n/a | yes |
| <a name="input_on_host_maintenance"></a> [on\_host\_maintenance](#input\_on\_host\_maintenance) | The maintenance policy for the runner | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy all resources into | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy all resources into | `string` | n/a | yes |
| <a name="input_runner_dir"></a> [runner\_dir](#input\_runner\_dir) | The directory to run the runner in | `string` | n/a | yes |
| <a name="input_runner_user"></a> [runner\_user](#input\_runner\_user) | The user to run the runner as | `string` | n/a | yes |
| <a name="input_runner_version"></a> [runner\_version](#input\_runner\_version) | GitHub Runner version to install | `string` | n/a | yes |
| <a name="input_subnet_name"></a> [subnet\_name](#input\_subnet\_name) | The subnet to deploy runner instances into | `string` | n/a | yes |
| <a name="input_use_spot_vms"></a> [use\_spot\_vms](#input\_use\_spot\_vms) | Whether to use spot VMs for the runner | `bool` | n/a | yes |
| <a name="input_vpc_name"></a> [vpc\_name](#input\_vpc\_name) | The VPC to deploy runner instances into | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_instance_template_name"></a> [instance\_template\_name](#output\_instance\_template\_name) | n/a |
| <a name="output_runner_service_account_id"></a> [runner\_service\_account\_id](#output\_runner\_service\_account\_id) | n/a |
<!-- END_TF_DOCS -->