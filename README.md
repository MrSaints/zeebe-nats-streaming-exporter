# zeebe-nats-streaming-exporter

_A work in progress._

An exporter for [Zeebe](https://github.com/zeebe-io/zeebe) to publish records [NATS Streaming Server (STAN)](https://nats.io/).

**WARNING: This is not ready for production. Please use at your own risk.**


---


## Getting Started

As this project is still in development, there are no pre-built JARs.

To use this with Zeebe:

1. `mvn package`
2. Copy the self-contained JAR (`target/zeebe-nats-streaming-exporter-$VERSION-jar-with-dependencies.jar`) to Zeebe `lib`
3. Register the exporter by adding the following to `zeebe.cfg.toml`:

```
[[exporters]]
id = "nats-streaming"
className = "com.fyianlai.zeebe.exporter.nats.NatsStreamingExporter"

  [exporters.args]
  # The available configuration options, and their defaults are below.
  # Uncomment, and configure accordingly.
  serverUrl = "nats://localhost"
  clusterId = "zeebe"
  clientIdPrefix = "zeebe-exporter-"
  channel = "zeebe"
  maxPubAcksInFlight = 10000
  format = "proto"
```


## Exported Data

By default, the data that is exported will be in `proto` format.

The format can be changed to `json` by configuring the following in `zeebe.cfg.toml`:

```diff
[exporters.args]
# ... other configuration ...
-      format = "proto"
+      format = "json"
```

The protobuf schema definition can be found in [`zeebe-exporter-protobuf`](https://github.com/zeebe-io/zeebe-exporter-protobuf/).


## Caveats

- No auto-reconnect if NATS Streaming Server goes down or if there are any intermittent connectivity problems
- If the number of published messages that goes unacknowledged reaches `maxPubAcksInFlight`, no further messages will be published, and the last exported record cursor will not be updated
- There is currently no plan for filtering messages before they are published, as this exporter solution is intended as a quick fire-and-forget
- For a similar reason to the previous point, there is also currently no plan for routing messages based on their `ValueType`
