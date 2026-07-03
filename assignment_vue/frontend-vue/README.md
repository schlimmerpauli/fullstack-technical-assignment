# Frontend Package Notes

This README is scoped to the Vue frontend package. For the full submission overview, start with the repository root `README.md`.

## Prerequisites

- Node.js 18 or higher

## Running the Development Server

From the `assignment_vue/frontend-vue/` directory:

```bash
npm install
npm run dev
```

The application is available at `http://localhost:5173`.

## Backend Connection

The frontend expects the backend to be running on `http://localhost:8080` during local development.

For separate frontend/backend deployments or a different API origin, set:

```bash
VITE_API_BASE_URL=http://localhost:8080
```

An example env file is available at `.env.example`.

## Validation Commands

From the `assignment_vue/frontend-vue/` directory:

```bash
npm test
npm run typecheck
npm run build
```

## Tech Stack

- Vue 3 with the Composition API and `<script setup>`
- TypeScript
- Vite
- Tailwind CSS

