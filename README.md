`k8split`
---

A CLI for splitting multidocument yaml files into discrete documents. Useful when converting k8s manifests to be fed into `kustomize`.

## Details

Documents are named by their "kind" and a numerical postfix. For example, a document with two `Deployment` objects would produce the files `deployment.yaml` and `deployment_1.yaml`.

## Installation

```console
go get -u github.com/brendanjryan/k8split
```

## Usage

Default:
```console
k8split file.yaml
``` 

To specify an output folder:

```console
k8split -o /split file.yaml
```


## Example 

```console
⇒  k8split config.yaml
2019/09/23 15:45:55 splitting config.yaml...
2019/09/23 15:45:55 split file into 8 chunks
2019/09/23 15:45:55 Writing file: namespace.yaml
2019/09/23 15:45:55 Writing file: deployment.yaml
2019/09/23 15:45:55 Writing file: ingress.yaml
2019/09/23 15:45:55 Writing file: service.yaml
2019/09/23 15:45:55 Writing file: cluster_role.yaml
2019/09/23 15:45:55 Writing file: service_account.yaml
2019/09/23 15:45:55 Writing file: cluster_role_binding.yaml
2019/09/23 15:45:55 Writing file: deployment_1.yaml
```
