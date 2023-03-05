---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_staging_vclconf Data Source - transparentedge"
subcategory: ""
description: |-
  
---

# transparentedge_staging_vclconf (Data Source)



## Example Usage

```terraform
data "transparentedge_staging_vclconf" "vclconfig" {}

output "staging_vcl_config" {
  value = data.transparentedge_staging_vclconf.vclconfig
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `company` (Number) Company ID that owns this Staging VCL config
- `id` (Number) ID of the Staging VCL Config
- `productiondate` (String) Date when the configuration was fully applied in the CDN
- `uploaddate` (String) Date when the configuration was uploaded
- `user` (String) User that created the configuration
- `vclcode` (String) Verbatim of the VCL code

