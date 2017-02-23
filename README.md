# Heimdall
[Heimdall][] is a monitoring service built on some of the same design
principles as other distributed monitoring tools like [Sensu][] but with
a strong emphasis on security and performance.

Heimdall accomplishes this by providing support for message authentication,
a strict security model and a secure-by-default design approach. In addition
to this, it is developed on Go, is automatically tested to ensure compatibility
and quality and has been built to enable easy deployment anywhere.

## Features
 - **Reliable Communication** over configurable transport layers, everything from
   [Redis][] to [NATS][], with the option to use multiple transports if required.
 - **Simple API** making it trivially easy to integrate with Heimdall, whether you're
   building a dashboard or an automated remediation tool.
 - **Easy to Deploy** whether you're using Docker or VMs thanks to its single executable.
 - **Aggregated Data** to simplify reporting, enable complex grouping and drive
   realtime dashboards.
 - **Client Driven Configuration**, since what you're monitoring matters more than
   your monitoring tool.

## Background
Monitoring of large scale distributed systems tends to be rather challenging to get
right. We don't believe it should be, so we've built Heimdall from the ground up with
the goal of making it trivially easy to monitor everything from containers to bare
metal servers, on any platform, in any network environment.

To accomplish this, Heimdall builds on the shoulders of proven, high performance,
communication services like [NATS][] and [Redis][] to ensure reliable and timely
delivery of information from your services to the Heimdall servers for processing.

With Heimdall, we also believe you shouldn't need to compromise on security in order
to monitor your systems effectively. Our security model ensures that a compromised
machine is incapable of influencing the behaviour of other machines in your Heimdall
cluster; while the ability to run Heimdall behind your firewall prevents your information
from leaving your hands.

## Design
Heimdall is composed of four components: the client, transport, server and datastore.
Each has a dedicated and well-defined purpose, keeping their implementations simple and
easy to audit, while also reducing the scope for failure. Heimdall's job is to provide
monitoring, so we've delegated the job of safely transporting events and storing data to
proven, existing, services.

We also believe that you shouldn't be tied to our technology choices and that what works
for one team might not suit another. Heimdall's transports and datastores, as a result,
are built to enable easy switching between a number of different options. Choose the one
you're most familiar with, or which suits your requirements best and trust Heimdall to do
the rest.

## Configuration

### Server

```yaml
---
# Listen on :80 for Heimdall's REST API
listen: ":80"

transports:
    # Define a NATS transport for the server to listen to
    # using the heimdall queue prefix.
  - driver: nats
    url: nats://localhost:4222/heimdall
```

### Client

```yaml
---
# Declare the details of our client including its name
# and tags. Names should be unique and are used for grouping.
client:
  # You can use $ENV variables in your client's name if you wish
  name: $HOSTNAME
  tags:
    role: webserver

transports:
    # Tell the client to submit events on a NATS transport
    # using the heimdall queue prefix.
  - driver: nats
    # We can use $ENV variables in our URLs
    url: nats://$USER:$PASS@localhost:4222/heimdall

checks:
    # Define an Apache healthcheck
  - name: apache-healthcheck
    # You can use $ENV variables in your commands
    command: curl -D - http://localhost:$APACHE_PORT/healthz
    # Run the check every 30 seconds
    interval: 30s
    # Give the check 3 seconds to run before failing
    timeout: 3s
    aggregates:
      # Aggregate the check's results into the "webservers" group
      - webservers
```

[Heimdall]: https://en.wikipedia.org/wiki/Heimdall_(comics)
[NATS]: http://nats.io/
[Redis]: https://redis.io/
[Sensu]: https://sensuapp.org/