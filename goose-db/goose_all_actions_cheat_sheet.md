
# 🐤 Goose CLI – All Actions and SQL Capabilities

> A complete cheat sheet for SQL-based migrations using Goose (CLI)

---

## 🔧 CLI Commands

| Command                          | Description |
|----------------------------------|-------------|
| `goose create NAME sql`          | Creates a new timestamped `.sql` migration file |
| `goose up`                       | Applies all pending `-- +goose Up` migrations |
| `goose up-to VERSION`           | Migrates up to a specific version |
| `goose down`                     | Rolls back the most recent migration |
| `goose down-to VERSION`         | Rolls back down to a specific version |
| `goose redo`                     | Rolls back and re-applies the most recent migration |
| `goose reset`                    | Rolls back all migrations (resets DB to version 0) |
| `goose status`                   | Shows migration status for each file |
| `goose version`                  | Shows the current version in the database |
| `goose fix`                      | Renames migrations to use consistent version prefixes |
| `goose validate`                | Validates migration consistency (Goose v3+) |

---

## 📁 Migration File Structure

```sql
-- +goose Up
-- +goose StatementBegin
-- SQL to apply during `up`
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQL to run during `down`
-- +goose StatementEnd
```

---

## ✅ SQL Operations You Can Perform

| SQL Operation       | Example |
|---------------------|---------|
| Create table        | `CREATE TABLE users (...);` |
| Drop table          | `DROP TABLE users;` |
| Alter table         | `ALTER TABLE users ADD COLUMN email TEXT;` |
| Insert default data | `INSERT INTO roles VALUES (1, 'admin');` |
| Update values       | `UPDATE users SET active = true;` |
| Delete rows         | `DELETE FROM logs WHERE id < 1000;` |
| Create index        | `CREATE INDEX ON users(name);` |
| Add constraint      | `ALTER TABLE orders ADD CONSTRAINT ...;` |
| Triggers/functions  | `CREATE FUNCTION notify_user() RETURNS trigger ...;` |

---

## 🌐 Example DSN for PostgreSQL

```bash
goose -dir goose-db postgres "host=localhost user=tasks_user password=task_9360 dbname=Tasks sslmode=disable" up
```

---

## 🧠 Pro Tips

- Always include a proper `-- +goose Down` block for rollbacks.
- Use `goose status` to confirm what's applied or pending.
- Structure your migrations per feature or table for clarity.

---
