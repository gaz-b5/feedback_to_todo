
# Feedback to Todo

## Overview

**Feedback to Todo** is a full-stack application designed to streamline the process of converting user feedback into actionable tasks. Built with a Go backend and a Next.js/React frontend, it provides a collaborative platform for managing projects, tasks, emails, and user feedback. The system leverages modern web technologies and database integrations to ensure scalability, reliability, and ease of use.

---

## Features

- **User Authentication:** Secure login, registration, and password reset flows.
- **Project Management:** Create, edit, and manage multiple projects.
- **Task Management:** Assign, track, and update tasks linked to projects and feedback.
- **Email Integration:** Process and associate emails with tasks and projects.
- **Dashboard:** Centralized view for managing all activities.
- **Database Support:** Uses SQLite for persistent storage (can be swapped for other DBs).
- **LLM Integration:** (Planned/Optional) Integrate with language models for feedback analysis.

---

## Project Structure

```
feedback_to_todo/
├── frontend/        # Next.js/React frontend
│   ├── app/         # Main app pages and API routes
│   ├── components/  # Reusable UI components
│   ├── data/        # Data schemas and seeds
│   ├── hooks/       # Custom React hooks
│   ├── lib/         # Utility functions
│   ├── public/      # Static assets
│   └── ...
├── server/          # Go backend
│   ├── controllers/ # Business logic (users, etc.)
│   ├── data_point/  # Data point logic
│   ├── formfunctions/ # Form handling
│   ├── llm_functions/ # LLM integration
│   ├── migrations/  # DB migrations
│   ├── models/      # ML models (optional)
│   ├── pb_data/     # Database files
│   ├── process_emails/ # Email processing
│   ├── qdrant_api/  # Vector DB integration
│   ├── routes/      # API routes
│   └── main.go      # Entry point
└── README.md        # Project documentation
```

---

## Prerequisites

- **Node.js** (v18+ recommended)
- **pnpm** (or npm/yarn)
- **Go** (v1.23 recommended)
- **Docker** (optional, for containerized deployment)

---

## Setup & Running Locally

### 1. Clone the Repository

```bash
git clone https://github.com/gaz-b5/feedback_to_todo.git
cd feedback_to_todo
```

### 2. Backend (Go Server)

#### a. Install Go dependencies

```bash
cd server
go mod tidy
```

#### b. Run the server

```bash
go run main.go serve
```

- The server will start on the configured port (default: 8090).
- Database files are stored in `server/pb_data/`.
- Environment variables can be set in `server/.env`.

#### c. (Optional) Run with Docker

```bash
cd server
docker build -t feedback-to-todo-server .
docker run -p 8080:8080 feedback-to-todo-server
```

### 3. Frontend (Next.js)

#### a. Install dependencies

```bash
cd frontend
pnpm install
# or
npm install
```

#### b. Run the development server

```bash
pnpm run dev
# or
npm run dev
```

- The frontend will start on [http://localhost:3000](http://localhost:3000).
- Environment variables can be set in `frontend/.env`.

---

## Usage

1. **Register/Login:** Create an account or log in.
2. **Create Projects:** Add new projects from the dashboard.
3. **Add Tasks:** Link feedback or emails to tasks within projects.
4. **Manage Members:** Add or edit project members.
5. **View & Process Emails:** See and process emails related to tasks.
6. **Dashboard:** Use the dashboard for an overview and quick actions.

---

## Database

- **Backend:** Uses SQLite databases stored in `server/pb_data/`.
- **Frontend:** No direct database access; communicates via API routes.
- **Migrations:** Managed via Go migration files in `server/migrations/`.

---

## Environment Variables

- **Backend:**
	- `PORT`: Server port (default: 8090)
	- `GROQ_API_KEY`: GROQ API key for llm
    - `QDRANT_DB_KEY`: QDRANT DB key for vector store
- **Frontend:**
	- API endpoint URLs
	- Auth keys (if needed)

---

## Deployment

- **Docker:** Use provided Dockerfile for backend containerization.
- **Frontend:** Deploy via Vercel, Netlify, or any Node.js host.
- **Environment:** Set production environment variables as needed.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request

---

## License

This project is licensed under the MIT License.

---

## Contact & Support

- GitHub Issues: [gaz-b5/feedback_to_todo](https://github.com/gaz-b5/feedback_to_todo/issues)
- Maintainer: gaz-b5

---

## Acknowledgements

- Built with Go, Next.js, Pocketbase, and SQLite
