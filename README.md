# secureapi-boatswain
Service responsible for running tests and sending results back to secureapi-web.

## generate binaries

```
make build
```

## container

```
make container
```

## kubernetes


### helm install for local development

To install the container built with `make container`

```
helm install --set image.repository=dev_kube/saucelabs/chef-scheduler --set image.tag=dev chart/chef-scheduler
```
