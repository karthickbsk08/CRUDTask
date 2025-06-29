
# 🐤 How to Install Goose on Ubuntu 22.04

Goose is a database migration tool. These steps show how to install the **Goose CLI** for managing SQL migrations.

---

## ✅ Prerequisites

- Go 1.18+ must be installed
```bash
go version
```

If not installed:
```bash
sudo apt update
sudo apt install golang -y
```

---

## 🛠️ Step-by-Step Goose Installation

### 1. Install Goose via `go install`

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

> This installs the `goose` binary in your `$HOME/go/bin` directory.

---

### 2. Add Goose to PATH

Edit your shell config file (`~/.bashrc` or `~/.zshrc`):

```bash
export PATH="$PATH:$HOME/go/bin"
```

Then apply the change:

```bash
source ~/.bashrc   # or ~/.zshrc
```

---

### 3. Verify Installation

```bash
goose -h
```

You should see the Goose CLI help output.

---

## 🧪 Test Example

Try running:

```bash
goose -version
```

Expected output:

```
goose version: <latest-version>
```

---

## 📦 Where It Installs

| Location             | Purpose                         |
|----------------------|----------------------------------|
| `$HOME/go/bin/goose` | The Goose CLI binary             |
| `$HOME/go/`          | Your Go workspace (default path) |

---

## 🧰 Bonus: Install Globally (Optional)

If you want to move Goose to a global bin directory:

```bash
sudo mv $HOME/go/bin/goose /usr/local/bin/
```

Now you can run `goose` from anywhere.

---
