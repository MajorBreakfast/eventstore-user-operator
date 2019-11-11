# Event Store User Operator

The Event Store user operator creates Event Store users according to Kubernetes resources in the cluster.

```
apiVersion: josefbrandl.com/v1
kind: EventStoreUser
metadata:
  name: my-app
spec:
  eventStore: my-eventstore
  groups:
    - foo
    - bar
```

This will create a `my-app-eventstore-user` secret with a `username` and `password` field in the same namespace.
