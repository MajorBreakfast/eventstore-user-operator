你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# Event Store User Operator

The Event Store user operator creates Event Store users according to `EventStoreUser` resources.

```yaml
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
  - Create a new user on the Event Store specified by `spec.eventStore` and assign it the user groups specified in `spec.groups`
  - Create a `my-app-eventstore-user` secret with a `username` and `password` field and put it in the same namespace
- If the `spec.groups` array on the `EventStoreUser` resource is updated on, the operator will update the user's groups on the Event Store.
- If the `EventStoreUser` resource is deleted, the operator will delete the user from the Event Store
- If the secret is deleted, the operator will reset the password of the user on the Event Store and create a new secret

Note: The `spec.eventStore` field should be considered immutable. Changing it, won't properly delete the user from the old Event Store.

The Event Store user operator is configured via the `/etc/eventstore-user-operator/config/config.yaml`:

```yaml
eventStores:
  - name: my-eventstore
    url: http://my-eventstore.my-eventstore.svc.cluster.local
```

- The `name` field defines the name under which `EventStoreUser` resources can reference the Event Store via their `spec.eventStore` field.
- The `url` should point to the HTTP endpoint of the Event Store (usually port 2113 on the Event Store container)

For each Event Store, the operator expects a `username` and `password` file under `/etc/eventstore-user-operator/eventstore-credentials/<eventstore-name>/`.
