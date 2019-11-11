# Event Store User Operator

The Event Store user operator creates Event Store users according to Kubernetes resources in the cluster.

```
apiVersion: josefbrandl.com/v1
kind: EventStoreUser
metadata:
  name: example-${i}
spec:
  eventStore: my-eventstore
  groups:
    - foo
    - bar
```
