# Database Schema Specification

This document details the database schema, field descriptions, and GORM ORM models for HY-Board.

---

## 💾 Database Engines & Compatibility

HY-Board supports three database dialects via GORM:
1. **SQLite 3** (Default, zero configuration, stored locally in `hy-board.db`)
2. **MySQL 8.0+ / MariaDB** (Production-ready, handles massive concurrency)
3. **PostgreSQL 14+** (Highly performant, stable enterprise-level DB)

The dialect and connection parameters are dynamically resolved at runtime from `config.yaml`.

---

## 📊 Entity Relationship & Table Schemas

### 1. `users` Table
Stores client login credentials, subscription status, limits, and proxy credentials.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique user identifier |
| `created_at` | `timestamp` | | Record creation timestamp |
| `updated_at` | `timestamp` | | Record last update timestamp |
| `deleted_at` | `timestamp` | Index | Soft delete timestamp |
| `email` | `string(255)` | Unique Index | User login email |
| `password` | `string(255)` | | Bcrypt hashed password |
| `v2ray_uuid` | `string(36)` | Unique Index | VMess/VLESS user ID (UUIDv4) |
| `trojan_password`| `string(64)` | | Trojan connection password |
| `speed_limit` | `uint32` | | Maximum bandwidth allowed in Mbps (0 = unlimited) |
| `device_limit` | `uint32` | | Maximum simultaneous device connections (0 = unlimited) |
| `total_traffic` | `uint64` | | Maximum allowed traffic in Bytes (0 = unlimited) |
| `used_traffic` | `uint64` | | Accumulated upload + download traffic in Bytes |
| `expired_at` | `timestamp` | | Subscription expiration date |
| `status` | `int8` | | Account status (1 = Active, 0 = Disabled) |
| `balance` | `double` | | User account balance |

---

### 2. `nodes` Table
Stores proxy node configurations and connection settings.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique node identifier |
| `created_at` | `timestamp` | | Record creation timestamp |
| `updated_at` | `timestamp` | | Record last update timestamp |
| `deleted_at` | `timestamp` | Index | Soft delete timestamp |
| `name` | `string(100)` | | User-facing node name |
| `type` | `string(20)` | | Protocol type: `V2ray`, `Vless`, `Trojan`, `Shadowsocks` |
| `address` | `string(255)` | | Server connection domain or IP |
| `port` | `uint16` | | Node port number |
| `traffic_rate` | `float32` | | Billing multiplier for user traffic (e.g. 1.0, 1.5, 0.5) |
| `settings` | `text` | | JSON string representing protocol-specific details |
| `show` | `boolean` | | Display status in user dashboard |

---

### 3. `traffic_logs` Table
A time-series table logging traffic usage on a per-user, per-node basis.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint64` | PK, Auto Increment | Log unique identifier |
| `user_id` | `uint` | Index | Foreign Key pointing to `users(id)` |
| `node_id` | `uint` | Index | Foreign Key pointing to `nodes(id)` |
| `up` | `uint64` | | Upload traffic generated in Bytes |
| `down` | `uint64` | | Download traffic generated in Bytes |
| `created_at` | `timestamp` | Index | Time the log entry was flushed |

---

### 4. `announcements` Table
Stores admin-published system announcements.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique announcement ID |
| `created_at` | `timestamp` | | Creation timestamp |
| `updated_at` | `timestamp` | | Last update timestamp |
| `deleted_at` | `timestamp` | Index | Soft delete timestamp |
| `title` | `string(255)` | | Announcement title |
| `content` | `text` | | Markdown content description |
| `show` | `boolean` | | Visibility status |

---

### 5. `knowledges` Table
Stores markdown articles for the user-facing knowledge base.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique article ID |
| `created_at` | `timestamp` | | Creation timestamp |
| `updated_at` | `timestamp` | | Last update timestamp |
| `deleted_at` | `timestamp` | Index | Soft delete timestamp |
| `title` | `string(255)` | | Article title |
| `content` | `text` | | Markdown content description |
| `show` | `boolean` | | Visibility status |
| `sort` | `integer` | | Ordering priority weight |

---

### 6. `tickets` Table
Stores user-submitted support tickets.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique ticket ID |
| `created_at` | `timestamp` | | Creation timestamp |
| `updated_at` | `timestamp` | | Last update timestamp |
| `deleted_at` | `timestamp` | Index | Soft delete timestamp |
| `user_id` | `uint` | Index | Proposer user ID pointing to `users(id)` |
| `title` | `string(255)` | | Ticket subject/title |
| `status` | `string(20)` | | Status (`open`, `replied`, `closed`) |

---

### 7. `ticket_messages` Table
Stores conversation messages inside a support ticket thread.

| Field Name | Data Type | Key / Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `uint` | PK, Auto Increment | Unique message ID |
| `created_at` | `timestamp` | | Creation timestamp |
| `ticket_id` | `uint` | Index | Belongs to ticket ID pointing to `tickets(id)` |
| `user_id` | `uint` | Index | Message sender user ID pointing to `users(id)` |
| `message` | `text` | | Message body |
| `is_admin` | `boolean` | | Sent by admin/staff status |

---

## 🛠️ Auto-Migration and Initialization

During backend startup via the Cobra `dev` or `prod` commands, GORM performs safe auto-migrations:
* Missing columns are added.
* Indexes are created.
* Existing data is preserved.
* A default admin user can be created running: `./hy-board admin init --email admin@example.com --password supersecure`.
