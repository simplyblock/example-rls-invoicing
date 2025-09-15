# Multi‑Tenancy with Row‑Level Security (RLS) in Postgres – Demo App

This repository is a small demo/example application showing how to implement multi‑tenancy using PostgreSQL Row‑Level
Security (RLS) from a Go backend using Goyave and pgx.

It demonstrates:

- A per‑tenant data model
- How to enable and enforce Postgres RLS policies
- How to connect and propagate tenant context in the application layer (Go + Goyave + pgx)

## Tech Stack

- Go 1.25
- [Goyave](https://goyave.dev/) (web framework)
- [pgx](https://github.com/jackc/pgx) (PostgreSQL driver)
- PostgreSQL 13+ (recommended)

## Repository Structure

- `database/schema.sql` – schema objects (tables, indexes, roles as needed)
- `database/demo-data.sql` – sample/demo data
- `database/rls.sql` – RLS setup (enabling RLS and defining policies)

Run the SQL in this order:

1) `schema.sql`
2) `demo-data.sql`
3) `rls.sql`

This order ensures your data loads before policies are enforced.

## Prerequisites

- Go 1.25 installed
- PostgreSQL server accessible
- psql CLI (or an equivalent tool)
- A database and a user with privileges to create objects and policies

## Database Setup

You need a PostgreSQL connection string. Example:

- For local dev:
    - Create a database: `createdb invoicing`
- Example connection string:
    - `postgres://postgres:your_password@localhost:5432/invoicing`

Load the SQL files in order:

Using psql:

```shell script
# 1) Schema
psql "$DATABASE_URL" -f database/schema.sql

# 2) Demo data
psql "$DATABASE_URL" -f database/demo-data.sql

# 3) RLS policies
psql "$DATABASE_URL" -f database/rls.sql
```

Alternatively, specify host, db, user explicitly:

```shell script
psql -h localhost -p 5432 -U postgres -d invoicing -f database/schema.sql
psql -h localhost -p 5432 -U postgres -d invoicing -f database/demo-data.sql
psql -h localhost -p 5432 -U postgres -d invoicing -f database/rls.sql
```

Using Docker (optional):

```shell script
docker run --name pg-rls-demo -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=invoicing -p 5432:5432 -d postgres:16

# Wait a few seconds for Postgres to become ready, then:
docker cp database/schema.sql pg-rls-demo:/schema.sql
docker cp database/demo-data.sql pg-rls-demo:/demo-data.sql
docker cp database/rls.sql pg-rls-demo:/rls.sql

docker exec -it pg-rls-demo psql -U postgres -d invoicing -f /schema.sql
docker exec -it pg-rls-demo psql -U postgres -d invoicing -f /demo-data.sql
docker exec -it pg-rls-demo psql -U postgres -d invoicing -f /rls.sql
```

Anyway, no matter how you arrived here, at this point, ensure you connect to the database with superuser privileges if
you want to create new records in the RLS-enabled tables.

## Build and Run

You must provide a `DATABASE_URL` environment variable when starting the application.

Build:

```shell script
go build ./...
```

Run (common patterns):

```shell script
# If main package is at repository root:
DATABASE_URL="postgres://postgres:your_password@localhost:5432/invoicing" go run .
```

If you use a binary:

```shell script
go build -o rls-demo .
DATABASE_URL="postgres://postgres:your_password@localhost:5432/invoicing" ./rls-demo
```

## How RLS Fits the Multi‑Tenancy Model

High‑level approach:

- Tables contain a tenant identifier column (for example, tenant_id) or have another resolvable reference that can be
  queries and enforced at query time.
- RLS is enabled on tenant‑scoped tables.
- Policies restrict SELECT/INSERT/UPDATE/DELETE to rows whose tenant_id matches the current tenant context.
- The app sets tenant context per request, typically via:
    - A session variable (e.g., `SET LOCAL app.current_tenant = '...'`) or
    - `SET ROLE role`/`SET LOCAL` GUCs that policies reference,
    - Using a connection‑pool hook to set the tenant on checkout,
    - or a combination of the above.

The `rls.sql` file contains the RLS policies used by this demo. Once enabled, queries that don’t satisfy the policy will
return no rows or fail on write.

Tip: Load demo data before enabling RLS (as shown) to avoid blocked inserts while policies are active.

## Troubleshooting

- Connection errors:
    - Verify DATABASE_URL, host, port, user, password, db name
    - Ensure Postgres is running and reachable
- Permission errors:
    - Ensure your DB user can create tables/policies or run the SQL as a superuser during setup
- No data returned:
    - RLS may be filtering results. Confirm tenant context is set and matches data tenant_id
- Data inserts fail:
    - RLS may require the insert row to match the current tenant_id; confirm policies and the values being inserted

## Notes

- The demo uses pgx; no additional DSN parameters are required beyond your environment needs.
- The Goyave framework handles routing and request lifecycle. Tenant context is attached per request and propagated to
  DB calls.

Feel free to explore the `database/` SQL files to understand the schema and policy rules, then inspect the Go code to
see how tenant context is applied on each request.

## Information

This example application is provided by [simplyblock](https://www.simplyblock.io) as is and does not intend to be
complete or production-ready. 
