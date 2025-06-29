# 🐤 Goose CLI Cheat Sheet for SQL Migrations

_Last updated: 2025-06-29 11:24:32_

---

## 🔨 1. Create a New Migration File

```bash
goose -dir goose-db create create_users_table sql
```

➡️ Creates a new migration file in `goose-db/` like:  
`goose-db/20250629124500_create_users_table.sql`

---

## 🚀 2. Apply Migrations (Up)

```bash
goose -dir goose-db postgres "<dsn>" up
```

Runs all **pending `-- +goose Up` migrations** in order.

---

## 🔙 3. Rollback the Last Migration (Down)

```bash
goose -dir goose-db postgres "<dsn>" down
```

Rolls back the most recent migration using `-- +goose Down`.

---

## 🎯 4. Migrate Up To a Specific Version

```bash
goose -dir goose-db postgres "<dsn>" up-to <version>
```

Run migrations up **to** a specific version number.

---

## 🔁 5. Rollback To a Specific Version

```bash
goose -dir goose-db postgres "<dsn>" down-to <version>
```

Rollback migrations **down to** a specific version.

---

## 📋 6. Check Current Migration Status

```bash
goose -dir goose-db postgres "<dsn>" status
```

Lists all migrations with "Applied" or "Pending" state.

---

## 🔢 7. Show Current DB Version

```bash
goose -dir goose-db postgres "<dsn>" version
```

Displays the **latest migration version** applied.

---

## 🧾 Sample Migration File

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
```

---

## 📘 Example DSN (PostgreSQL)

```bash
goose -dir goose-db postgres "host=localhost user=tasks_user password=task_9360 dbname=Tasks sslmode=disable" up
```

---
