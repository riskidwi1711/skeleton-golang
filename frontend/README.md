# ITMS Frontend

Vue 3 + Vite admin UI starter untuk ITMS.

## Current pages
- `/dashboard` (tenant summary dari API)
- `/tickets` (UI draft)
- `/assets` (UI draft)
- `/onboarding` (form onboarding tenant)

## Run local
```bash
cd frontend
npm install
npm run dev
```

## Build
```bash
npm run build
npm run preview
```

## API target
Default API base:
- `/api` (same-origin proxy from `itms.riski-labs.site` to gateway)

Override via `.env`:
```bash
VITE_API_BASE=https://itms-api.riski-labs.site
```
