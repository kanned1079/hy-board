# API Documentation

This document defines the REST API endpoints for HY-Board. The API is divided into three scopes:
1. **UniProxy API**: Used by the XrayR backend.
2. **Client API**: Used by users for dashboards, purchases, and subscriptions.
3. **Admin API**: Used by administrators to manage users, nodes, and configurations.

---

## 🔒 Authentication

### 1. XrayR Backend
XrayR authenticates via a static header set in `config.yaml`:
* **Header**: `Token`
* **Value**: `<secret_uniproxy_token>`

### 2. Client & Admin Dashboards
Dashboards use JWT bearer tokens:
* **Header**: `Authorization`
* **Value**: `Bearer <jwt_token>`

---

## 📡 1. UniProxy API (XrayR Integration)

All UniProxy endpoints are prefixed with `/api/v1/server/UniProxy`.

### 1.1 Get Node Configuration
Fetch the specific configuration details for a node.
* **Method**: `GET`
* **Path**: `/api/v1/server/UniProxy/config`
* **Query Parameters**:
  * `node_id` (required, int) - The unique ID of the node.
  * `node_type` (required, string) - The node protocol (e.g. `V2ray`, `Vless`, `Trojan`, `Shadowsocks`).
* **Request Header**: `Token: <token>`
* **Success Response (200 OK)**:
  ```json
  {
    "data": {
      "port": 443,
      "address": "node1.example.com",
      "transport": "ws",
      "path": "/v2ray",
      "enable_vless": true,
      "vless_flow": "xtls-rprx-vision"
    }
  }
  ```

### 1.2 Get Authorized Users
Retrieve the list of users permitted to connect to this node.
* **Method**: `GET`
* **Path**: `/api/v1/server/UniProxy/user`
* **Query Parameters**:
  * `node_id` (required, int)
  * `node_type` (required, string)
* **Request Header**: `Token: <token>`
* **Success Response (200 OK)**:
  ```json
  {
    "data": [
      {
        "id": 1,
        "uuid": "4f5e08b1-3ea9-4f76-8802-1249b6b802e3",
        "speed_limit": 100,
        "device_limit": 3
      }
    ]
  }
  ```

### 1.3 Push Node Status & Traffic Logs
Report real-time traffic usage and node status.
* **Method**: `POST`
* **Path**: `/api/v1/server/UniProxy/push`
* **Request Header**: `Token: <token>`
* **Request Body**:
  ```json
  {
    "node_id": 1,
    "traffic": [
      {
        "user_id": 1,
        "up": 1048576,
        "down": 4194304
      }
    ]
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "message": "success"
  }
  ```

---

## 👤 2. Client API

All client endpoints are prefixed with `/api/v1/client`.

### 2.1 User Login
* **Method**: `POST`
* **Path**: `/api/v1/client/auth/login`
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "mypassword"
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "token": "eyJhbGciOi...",
    "user": {
      "email": "user@example.com",
      "total_traffic": 107374182400,
      "used_traffic": 53687091200
    }
  }
  ```

### 2.2 Get Nodes List
* **Method**: `GET`
* **Path**: `/api/v1/client/nodes`
* **Headers**: `Authorization: Bearer <jwt>`
* **Success Response (200 OK)**:
  ```json
  {
    "nodes": [
      {
        "id": 1,
        "name": "US-01 Premium",
        "type": "Vless",
        "address": "node1.example.com",
        "port": 443,
        "traffic_rate": 1.0
      }
    ]
  }
  ```

### 2.3 Fetch Subscription Config
Returns raw subscription format (e.g. Clash YAML or Base64 string depending on User-Agent).
* **Method**: `GET`
* **Path**: `/api/v1/client/subscribe`
* **Query Parameters**:
  * `token` (required, string) - The user's specific subscription token.
  * `client` (optional, string) - Client identifier (e.g. `clash`, `shadowrocket`).
* **Success Response (200 OK)**: Raw subscription configuration contents matching the target client format.

---

## 🔑 3. Admin API

All admin endpoints are prefixed with `/api/v1/admin`.

### 3.1 Get User Accounts
* **Method**: `GET`
* **Path**: `/api/v1/admin/users`
* **Headers**: `Authorization: Bearer <jwt>`
* **Success Response (200 OK)**: Array of user records including billing and data usage.

### 3.2 Create Proxy Node
* **Method**: `POST`
* **Path**: `/api/v1/admin/nodes`
* **Headers**: `Authorization: Bearer <jwt>`
* **Request Body**:
  ```json
  {
    "name": "SG-01",
    "type": "Vless",
    "address": "sg1.example.com",
    "port": 443,
    "traffic_rate": 1.2,
    "settings": "{\"transport\": \"grpc\", \"path\": \"TuneTunnel\"}"
  }
  ```
* **Success Response (201 Created)**: Created Node object.

---

## ✉️ 4. Support Tickets API (Shared REST & WebSocket)

All ticket endpoints are prefixed with `/api/v1/tickets` or `/api/v1/ws`.

### 4.1 Get Tickets List
Fetch the list of tickets. Normal clients see only their own tickets, while admins see all tickets.
* **Method**: `GET`
* **Path**: `/api/v1/tickets`
* **Headers**: `Authorization: Bearer <jwt>`
* **Success Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "created_at": "2026-06-08T09:30:00Z",
      "updated_at": "2026-06-08T09:35:00Z",
      "user_id": 2,
      "user": {
        "id": 2,
        "email": "user@example.com"
      },
      "title": "Connection failure on SG node",
      "status": "open"
    }
  ]
  ```

### 4.2 Get Ticket Details
Fetch a specific ticket's details and conversation history.
* **Method**: `GET`
* **Path**: `/api/v1/tickets/:id`
* **Headers**: `Authorization: Bearer <jwt>`
* **Success Response (200 OK)**:
  ```json
  {
    "id": 1,
    "created_at": "2026-06-08T09:30:00Z",
    "updated_at": "2026-06-08T09:35:00Z",
    "user_id": 2,
    "user": {
      "id": 2,
      "email": "user@example.com"
    },
    "title": "Connection failure on SG node",
    "status": "replied",
    "messages": [
      {
        "id": 1,
        "created_at": "2026-06-08T09:30:00Z",
        "ticket_id": 1,
        "user_id": 2,
        "user": {
          "id": 2,
          "email": "user@example.com"
        },
        "message": "My connection keeps timing out.",
        "is_admin": false
      }
    ]
  }
  ```

### 4.3 Create New Ticket
Submit a new support ticket.
* **Method**: `POST`
* **Path**: `/api/v1/tickets`
* **Headers**: `Authorization: Bearer <jwt>`
* **Request Body**:
  ```json
  {
    "title": "Connection failure on SG node",
    "message": "My connection keeps timing out."
  }
  ```
* **Success Response (201 Created)**: The created Ticket object with preloaded user and messages.

### 4.4 Real-time Tickets WebSocket Channel
Connect to real-time ticket events (creation, reply, status closure).
* **Protocol**: WebSocket
* **Path**: `/api/v1/ws/tickets`
* **Query Parameters**:
  * `token` (required, string) - The user's JWT authorization token.
* **WebSocket Message Protocols (JSON)**:
  * **Actions (Client to Server)**:
    * `{ "action": "reply_ticket", "ticket_id": <int>, "message": "<string>" }`
    * `{ "action": "close_ticket", "ticket_id": <int> }`
  * **Events (Server to Client)**:
    * `{ "event": "tickets_list", "data": [...] }` (sent on initial connection)
    * `{ "event": "ticket_created", "data": { ... } }`
    * `{ "event": "ticket_updated", "data": { ... } }`

