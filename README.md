# HY-Board: High-Aesthetic XrayR Management Panel

HY-Board is a modern, high-aesthetic control panel designed to manage proxy server subscriptions and interface seamlessly with the **XrayR** backend. Built using Golang on the backend and Nuxt 4 (Vue 3) + Nuxt UI on the frontend, it delivers extreme performance and a stunning, responsive user experience.

---

## 🚀 Tech Stack

### Backend
* **Core**: Golang 1.22+
* **Routing**: Gin Web Framework
* **ORM**: GORM (supporting SQLite, MySQL 8+, and PostgreSQL)
* **CLI Management**: Cobra (`spf13/cobra`)
* **Caching & Queue**: Redis (optional, for high-concurrency traffic logging)
* **Configuration**: YAML (`config.yaml`)

### Frontend
* **Core**: Nuxt 4 / Vue 3
* **UI Framework**: Nuxt UI (Tailwind CSS, Headless UI, and Heroicons integrated)
* **Styling & Motion**: Custom CSS (Glassmorphism, non-linear ease-in-out animations)

---

## 📂 Project Structure

```text
hy-board/
├── README.md               # Project overview & system memory
├── DATABASE.md             # Database schema and ORM mappings
├── API.md                  # API documentation (UniProxy, Client, Admin)
├── backend/                # Go backend application
│   ├── config.yaml         # Configuration file (default: SQLite)
│   ├── main.go             # Entrypoint
│   ├── config/             # YAML configuration parser
│   ├── database/           # DB & Redis connection initializers
│   ├── models/             # GORM models & auto-migration
│   ├── middleware/         # Gin middlewares (Auth, CORS, Logger)
│   ├── cmd/                # Cobra CLI commands (dev, prod, admin)
│   └── routes/             # Gin route groups & controllers
└── frontend/               # Nuxt 4 frontend application
    ├── nuxt.config.ts      # Nuxt configuration (API devProxy enabled)
    ├── app.config.ts       # Global Nuxt UI theme & transition config
    ├── .env                # Runtime environment variables
    ├── app.vue             # Root component
    └── pages/              # Nuxt Pages (index, login, dashboard)
```

---

## 🤖 Contextual Memory (For AI Agents)

* **Architecture**: The backend serves as the control plane. Node servers run XrayR (data plane) and communicate with the backend's `UniProxy` REST API endpoints.
* **Database Driver Selection**: The database engine is configured in `config.yaml` (`database.type: sqlite | mysql | postgres`). GORM dynamically opens the database connection.
* **Routing Prefix**: All API calls from the client go through `/api/v1/...`. In development, Nuxt's Nitro dev server routes `/api/v1/**` requests to the Go backend using a proxy configured via `.env` and `nuxt.config.ts`.
* **Animations**: All transitions (page-to-page navigation, modal pop-ups) must use custom, non-linear timing functions (`cubic-bezier(0.25, 1, 0.5, 1)` or similar) to ensure premium aesthetic feedback.
