# Event Store User Operator

The Event Store user operator creates Event Store users according to `EventStoreUser` resources.

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

- After the `my-app` `EventStoreUser` resource is created, the operator will:
  - Create a new user on the Event Store and assign it the specified user groups
  - Create a `my-app-eventstore-user` secret with a `username` and `password` field and put it in the same namespace
- If the groups array on the `EventStoreUser` resource is updated on, the operator will update the user's groups on the Event Store.
- If the `EventStoreUser` resource is deleted, the operator will delete the user from the Event Store
- If the secret is deleted, the operator will reset the password of the corresponding user and create a new secret

Note: The `spec.eventStore` field should be considered immutable. Changing it, won't properly delete the user from the old Event Store.

The Event Store User operator is configured via the `/etc/eventstore-user-operator/config/config.yaml`:

```
eventStores:
  - name: my-eventstore
    url: http://my-eventstore.my-eventstore.svc.cluster.local
```

For each Event Store, the operator expects a `username` and `password` file under `/etc/eventstore-user-operator/eventstore-credentials/<eventstore-name>/`.
