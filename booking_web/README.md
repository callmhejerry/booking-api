# Booking Web

A small Flutter web client for the booking API. It can search slots, create
bookings, list bookings by email, fetch a booking by ID, cancel a booking, and
replace the slot schedule.

The page also includes an administrative `POST /reset` panel. It accepts a
`{"slots": [...]}` payload and requires confirmation because the API deletes
the existing slots (and potentially related bookings) before recreating them.

## Run locally

From this directory:

```sh
flutter pub get
flutter run -d web-server --web-hostname 0.0.0.0 --web-port 5000
```

Open `http://localhost:5000`, enter the API's base URL and `API_KEY`, then search for a provider. For a local API, the base URL is usually `http://localhost:8080`; for Render, use the API service's `https://...onrender.com` URL.

The API must allow the web app's origin in `CORS_ALLOWED_ORIGINS`. The default is `http://localhost:5000`. For a hosted frontend, set this variable on the API service to the frontend's exact origin (for example, `https://booking-web.onrender.com`). Multiple origins can be comma-separated.

## Deploy on Render

The repository's root `render.yaml` defines a free Render Static Site named
`booking-web`. Render installs the pinned Flutter SDK, builds this folder, and
publishes `build/web`. The build reuses the SDK directory on subsequent attempts.
Create or update the Render Blueprint from the repository
to deploy both the API and web app. The API CORS setting allows
`https://booking-web.onrender.com`; if Render assigns a different URL, update
the API service's `CORS_ALLOWED_ORIGINS` to match the exact frontend origin.

After both services are live, open the static site's URL, enter the API
service's `https://...onrender.com` URL and the API key from the API service's
environment settings, then search for a provider.

## Build for other static hosts

```sh
flutter build web --release
```

Publish the generated `build/web` directory with any static-file host. Add that
static site's origin to the API's `CORS_ALLOWED_ORIGINS`.

The API key is entered by the user and held only in page memory; do not hard-code a private key into a public web build. A key distributed to browser users should be considered public.
