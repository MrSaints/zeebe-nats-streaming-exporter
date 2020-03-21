# sql

_A work in progress._

A Postgresql importer for Zeebe NATS Streaming exporter.

**WARNING: This is not ready for production. Please use at your own risk.**


---


## Getting Started

To set-up the database, you can use [Soda CLI](https://gobuffalo.io/en/docs/db/toolbox) (see commands below) or you can apply the `migrations/schema.sql` to your Postgresql database.

For the former, you need to add your Postgresql environment in `database.yml`, and run the commands in this directory (where this `README.md` lives).

```
# Create Database
soda create -a

# Run Migrations
soda migrate up

# Destroy Database
soda drop -a
```

As this project is still in development, there are no packaged binary or Docker image for the importer.

To run the importer, `go run *.go`.

Various credentials are currently hardcoded, but there will be work to make them configurable.


## TODO

- [ ] Build proper CLI with configurations
- [ ] Handle `VariableRecord`
- [ ] Handle `IncidentRecord`
- [ ] Handle `MessageRecord`
- [ ] Handle `MessageSubscriptionRecord`
