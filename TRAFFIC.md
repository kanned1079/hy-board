# Traffic Accounting Architecture

This document explains the complete, end-to-end flow of how user traffic is measured, reported, stored, and displayed in HY-Board.

---

## Overview

HY-Board uses **Xray-core's native gRPC StatsService** to collect precise, byte-level traffic statistics for each user. This replaces fragile log-scraping approaches and provides accurate upload/download counters per user, reset-on-read so no bytes are double-counted.

```
 User Device
     │ (proxy traffic)
     ▼
 ┌─────────────┐
 │  Xray-core  │  ← tracks per-user bytes via StatsService
 │  (running   │
 │  on node)   │
 └──────┬──────┘
        │ (gRPC query every N seconds)
        ▼
 ┌─────────────┐    HTTP POST /api/v1/server/UniProxy/push
 │    Daemon   │ ──────────────────────────────────────────► ┌─────────────┐
 │  (Go binary)│                                             │   Backend   │
 └─────────────┘                                             │  (Gin/GORM) │
                                                             └──────┬──────┘
                                                                    │
                                             ┌──────────────────────┤
                                             │                      │
                                             ▼                      ▼
                                      traffic_logs            users.used_traffic
                                       (time-series)          (running total)
```

---

## Step 1 — Xray Config: Enabling Per-User Stats

When the daemon generates the Xray config, it injects three extra blocks that activate the gRPC statistics API:

```json
{
  "stats": {},
  "api": {
    "tag": "api",
    "listen": "127.0.0.1:10085",
    "services": ["StatsService"]
  },
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true
    }
  }
}
```

Every inbound client is assigned an `email` tag in the format:

```
user_<userID>@hy-board.com
```

For example, user with ID `7` gets the email `user_7@hy-board.com`. Xray uses this email tag as the stats key to track that user's traffic independently.

**Source:** `daemon/main.go` → `fetchConfigFromApi()`

---

## Step 2 — Daemon: Polling Traffic from Xray

Inside the supervisor ticker loop (every `interval` seconds, default 15 s), the daemon runs:

```bash
xray api statsquery -s 127.0.0.1:10085 -pattern "user>>>" -reset
```

This calls Xray's `StatsService.QueryStats` gRPC method via the CLI wrapper. The `-reset` flag ensures each call reads only the **delta** since the last poll — no bytes are counted twice.

The raw JSON output looks like:

```json
{
  "stat": [
    { "name": "user>>>user_7@hy-board.com>>>traffic>>>uplink",   "value": 1048576 },
    { "name": "user>>>user_7@hy-board.com>>>traffic>>>downlink", "value": 5242880 }
  ]
}
```

The daemon parses the `name` field by splitting on `>>>`:

```
parts[0] = "user"
parts[1] = "user_7@hy-board.com"   ← extract userID from here
parts[2] = "traffic"
parts[3] = "uplink" | "downlink"
```

The numeric `value` is typed as `int64` bytes and cast to `uint64`.

**Source:** `daemon/main.go` → `queryTrafficFromXray()`

---

## Step 3 — Daemon: Reporting to Backend

After parsing, the daemon POSTs the traffic delta to the panel via:

```
POST /api/v1/server/UniProxy/push
Header: Token: <api_key>
```

```json
{
  "node_id": 1,
  "traffic": [
    { "user_id": 7, "up": 1048576, "down": 5242880 }
  ]
}
```

The request is fired in a separate goroutine so it does not block the supervisor tick loop.

**Source:** `daemon/main.go` → `reportTrafficToApi()`

---

## Step 4 — Backend: Storing Traffic

The `PushNodeData` handler processes the incoming payload inside a single database transaction:

```
POST /api/v1/server/UniProxy/push
```

For each entry in `traffic[]`:

1. **Increments `users.used_traffic`** atomically using a SQL expression:
   ```sql
   UPDATE users SET used_traffic = used_traffic + <up+down> WHERE id = <user_id>
   ```

2. **Inserts a new row into `traffic_logs`** with the raw up/down values, a `node_id`, and a `created_at` timestamp (for historical analysis).

**Source:** `backend/routes/uniproxy.go` → `PushNodeData()`

---

## Step 5 — GraphQL: Node Status from Traffic Logs

The `nodeType` GraphQL object exposes a computed `status` field that reads `traffic_logs` in real time to determine whether a node is actively serving users:

| Status | Condition |
|:---|:---|
| `"offline"` | No heartbeat from the daemon in the last 45 seconds (`nodes.updated_at`) |
| `"active"` | At least one `traffic_logs` row for this node with `up > 0 OR down > 0` within the last 60 seconds |
| `"idle"` | Daemon is alive (heartbeat OK) but no active user traffic in the last 60 seconds |

**Source:** `backend/routes/graphql_types.go` → `nodeType` → `"status"` field resolver

---

## Step 6 — Frontend: Live Dashboard Refresh

The user dashboard polls `userInfo { used_traffic, total_traffic, balance }` every **30 seconds** via a lightweight GraphQL query (separate from the full page load). This keeps the traffic display current without a full page reload.

A visual countdown timer (30 → 0 s) and a clickable manual refresh icon are shown on the Traffic Card. During a fetch the value fades to 50% opacity; a thin progress bar fills up as the countdown progresses.

**Source:** `frontend/pages/dashboard/index.vue` → `refreshTraffic()` + `startRefreshCycle()`

---

## Data Flow Summary

| Stage | Component | Interval | What happens |
|:---|:---|:---|:---|
| Traffic collection | **Xray-core** (StatsService) | Continuous | Tracks per-user bytes internally via `email` tag |
| Traffic query | **Daemon** (`xray api statsquery`) | Every N s (default 15 s) | Reads + resets per-user counters atomically |
| Traffic push | **Daemon → Backend** (`/push`) | Same interval | POSTs deltas to backend |
| Database write | **Backend** (GORM tx) | On each push | Increments `used_traffic`, appends `traffic_logs` row |
| Node status | **GraphQL** (real-time resolve) | On query | Checks `traffic_logs` within last 60 s |
| UI refresh | **Frontend** (Nuxt, `setInterval`) | Every 30 s | Fetches fresh `used_traffic` via GraphQL |

---

## Key Files

| File | Role |
|:---|:---|
| [`daemon/main.go`](daemon/main.go) | Xray lifecycle manager, traffic polling and push |
| [`backend/routes/uniproxy.go`](backend/routes/uniproxy.go) | `/push` HTTP endpoint, DB writes |
| [`backend/routes/graphql_types.go`](backend/routes/graphql_types.go) | `nodeType` status resolver reading `traffic_logs` |
| [`backend/models/models.go`](backend/models/models.go) | `TrafficLog` struct (GORM model) |
| [`frontend/pages/dashboard/index.vue`](frontend/pages/dashboard/index.vue) | Auto-refresh UI, countdown timer |

---

## Important Design Decisions

- **`-reset` flag is critical.** Without it, every poll would return the cumulative total since Xray started, causing massive double-counting. With `-reset`, each read is a delta covering exactly the last interval window.
- **`email` tag encodes the user ID**, not the actual email address, so the user ID can be extracted from the stat key without a database lookup at parse time.
- **Atomic SQL update** (`used_traffic + ?`) prevents race conditions when multiple daemon instances (multi-node setups) push simultaneously.
- **`traffic_logs` is append-only.** Old rows are never updated. This provides a full audit trail and enables future analytics (e.g., hourly/daily breakdowns per node).
