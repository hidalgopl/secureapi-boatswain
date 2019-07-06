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

## Test with curl
export AUTH_TOKEN=$(k get secret boatswain-web -o jsonpath="{.data.authKey}" | base64 --decode)
kubectl run curl --image=radial/busyboxplus:curl -i --tty
curl -X POST -H "X-SecureAPI-Secret-Key: $AUTH_TOKEN" -d '{"user_id": "fake", "tests": ["SEC#0001", "SEC#0002", "SEC#0006"], "url": "https://forumcm.net", "test_suite_id": "1"}' boatswain.default.svc.cluster.local:8080/tests/schedule