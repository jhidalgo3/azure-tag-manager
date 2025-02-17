# Azure Tag manager 

Tag manager is used to add or change tags on Azure resources on a massive scale. The tags are applied only when given rules are satisfied, for example prior existance of a certain tag, or the lack of it. 

[![Go Report Card](https://goreportcard.com/badge/github.com/jhidalgo3/azure-tag-manager)](https://goreportcard.com/report/github.com/jhidalgo3/azure-tag-manager)

[Manage resource groups with the Azure SDK for Go](https://learn.microsoft.com/en-us/azure/developer/go/manage-resource-groups?tabs=bash%2Cazure-portal)

## Download

You can download a binary for your architecture from the releases page. 

## Usage

To use it you need to create Azure service principal. Create service principal file:

```bash
az ad sp create-for-rbac
```
and export path to the authorizer:

```bash
export AZURE_SUBSCRIPTION_ID=<CHANGE_ME>
export AZURE_TENANT_ID=<CHANGE_ME>
export AZURE_CLIENT_ID=<CHANGE_ME>
export AZURE_CLIENT_SECRET=<CHANGE_ME>
```

## How it works ?

Azure Tag Manager needs a file with rules, each rule consists of conditions and actions. Given all the conditions, all actions are executed. The rules file can be supplied either in json or yaml. 

```json
{
  "dryrun": true,
  "rules":  [
    {
        "name": "Tag me stuff", 
        "conditions": [
            {"type": "tagEqual", "tag": "myTag", "value" : "example"},
            {"type": "tagExists", "tag": "myTagWhichExists"},
            {"type": "tagNotExists", "tag": "env"}
        ], 
        "actions": [
            {"type": "addTag", "tag": "newTag", "value": "value" },
            {"type": "addTag", "tag": "newTag2", "value": "value2" }
        ]
      }
    ]
}
```

Equivalent of the same file in YAML would look like:

```YAML
---
dryrun: true
rules:
- name: Tag me this
  conditions:
  - type: tagEqual
    tag: myTag
    value: example
  - type: tagExists
    tag: myTagWhichExists
  - type: tagNotExists
    tag: env
  - type: regionEqual
    region: westeurope
  - type: rgEqual
    resourceGroup: exampleRg
  actions:
  - type: addTag
    tag:  newTag
    value: newValue
  - type: addTag
    tag: newTag2
    value: someValue
```

The following types of conditions are accepted:

* `noTags` - checks if there are no tags set 
* `tagEqual` - checks if a `tag` has a `value` set 
* `tagNotEqual` - checks if a `tag` has a value set different than `value` 
* `tagExists` - checks if a tag with key `tag` exists
* `tagNotExists` - same as above but negative
* `regionEqual` - checks if resource is in key `region` (aka location in azure)
* `regionNotEqual` - same as above but negative
* `rgEqual` - match resource group in a key `resourceGroup`
* `rgNotEqual` - match not resource group
* `resEqual` - resource name equals `resource` 

The supported actions are:

* `addTag` - adds a tag with key `tag` and value `value`
* `delTag` - deletes a tag with key `tag`

When rewriting, the tool will first do a backup of old tags. It will be saved in a file in the current (run) directory. 

## Running 

Tagmanager accepts commands and flags: `tagmanager COMMAND [FLAGS`]. 
```
Usage:
  tagmanager [command]

Available Commands:
  check       Do sanity checks on a resource group (NOT FULLY IMPLEMENTED YET)
  help        Help about any command
  restore     Restore previous tags from a file backup
  retagrg     Retag resources in a rg based on tags on rgs
  rewrite     Rewrite tags based on rules from a file

Flags:
  -h, --help      help for tagmanager
  -v, --verbose   verbose output
```

Commands:

* `rewrite` - mode where tagmanager will retag the resources based on mapping given in a mapping file input (specified with `-m filepath` flag). If `--dry` flag is given, the tagging actions will not be executed

```
go run cmd/cli/main.go rewrite -m rules.yaml -v
```

* `restore` - restores tags backed up in a file, supplied by `-f filepath` flag

* `check` - (EXPERIMENTAL) does some basic sanity checks on the resource group given as `--rg` flag 

* `retagrg` - Takes tags form a given resource group (`--rg`) and applies them to all of the resources in the resource group. If any existing tags are already there, the new ones with be appended. Adding `--cleantags` will clean ALL the tags on resources before adding new ones. 

```
go run cmd/cli/main.go retagrg -r <RESOURCE_NAME> -v
```


## RUNNING



## MIGRATION AZURE GO SDK

[Upgrade to Azure GO SDK](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/MIGRATION_GUIDE.md)

## Thanks 

Original code from:

Author(s): Dariusz Dwornikowski (@tdi)
Made by Nordcloud in Poznań with ♥. 