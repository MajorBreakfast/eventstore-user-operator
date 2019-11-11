NUM_USERS=${1:-5}

MANIFESTS=$(

for i in $(seq 1 $NUM_USERS); do

cat << EOF
---
apiVersion: josefbrandl.com/v1
kind: EventStoreUser
metadata:
  name: example-${i}
spec:
  eventStore: my-eventstore
  groups:
    - foo
    - bar
EOF

done

)


echo "${MANIFESTS}" | kubectl apply -f -
