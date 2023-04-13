# dapr-im
dapr im caht

##  Dapr 消息 TTL 过期时间
```
from dapr.clients import DaprClient

with DaprClient() as d:
    req_data = {
        'order-number': '345'
    }
    # Create a typed message with content type and body
    resp = d.publish_event(
        pubsub_name='pubsub',
        topic='TOPIC_A',
        data=json.dumps(req_data),
        publish_metadata={'ttlInSeconds': '120'}
    )
    # Print the request
    print(req_data, flush=True)


[
  {
    "pubsubname": "pubsub",
    "topic": "newOrder",
    "route": "/orders",
    "metadata": {
      "rawPayload": "true",
    }
  }
]

```