# terraform-gcp-github-runner
<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_google"></a> [google](#requirement\_google) | ~> 5.16 |

## Providers

No providers.

## Modules

No modules.

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The GCP project ID to deploy all resources into | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy all resources into | `string` | `"europe-west1"` | no |
| <a name="input_zones"></a> [zones](#input\_zones) | The zones to deploy the runner instances into | `list(string)` | <pre>[<br>  "europe-west1-b",<br>  "europe-west1-c",<br>  "europe-west1-d"<br>]</pre> | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->