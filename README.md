# User Profile Application

A backend-focused web application built in Go with server-rendered frontend and a separate REST API layer.

## Features

- User Signup/Login (manual + Google OAuth)
- Session-based authentication
- Profile creation and update
- Role-based behavior (Google vs Manual users)
- Middleware-based auth validation
- REST API independent of frontend

---

## Tech Stack

- Backend: Go (net/http)
- Database: MySQL
- Frontend: Go Templates
- Auth: Google OAuth 2.0
- Deployment: Railway

---

## Architecture
Browser
↓
Go Templates (Page Handlers)
↓ HTTP
REST API (/api/*)
↓
Service Layer
↓
Repository Layer
↓
MySQL

The frontend communicates with the backend strictly via REST APIs.

--

## API Endpoints

Auth
- POST /api/login
- POST /api/signup
- GET /api/auth/google/callback

Profile
- GET /api/profile
- POST /api/profile/create
- POST /api/profile/update


## .env
```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=root
DB_NAME=yourdb

GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GOOGLE_USER_INFO=https://www.googleapis.com/oauth2/v2/userinfo

BASE_URL=http://localhost:8000
PORT=8000

```