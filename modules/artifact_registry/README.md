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
| [google_artifact_registry_repository.image_cache](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/artifact_registry_repository) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy the repository to. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy the repository to. | `string` | `"europe-west1"` | no |
| <a name="input_remote_repository_url"></a> [remote\_repository\_url](#input\_remote\_repository\_url) | The URL of the remote repository to clone. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_image_cache_url"></a> [image\_cache\_url](#output\_image\_cache\_url) | n/a |
<!-- END_TF_DOCS -->