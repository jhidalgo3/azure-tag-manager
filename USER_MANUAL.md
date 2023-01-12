# Retag all resources in resource group MAIN

```
go run cmd/cli/main.go retagrg -r MAIN -v --cleantags
```

# Rewrite tags from rules.yaml file

```
go run cmd/cli/main.go -v rewrite -m rules.yaml  --dry
```