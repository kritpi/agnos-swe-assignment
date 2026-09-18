# API Spec — Hospital Middleware

Base URL:
- through nginx: `http://localhost:${NGINX_PORT}` (default port `80`), which proxies every path to the Go service
- running the service directly (`make run`): `http://localhost:${PORT}` (default port `8080`)

The paths below are relative to the base URL. Content type: `application/json` for every request and response.

## Endpoints

| Method | Path | Purpose | Login required |
|---|---|---|---|
| GET | `/healthz` | Liveness check (outside `/api/v1`) | No |
| GET | `/docs/api` | Swagger UI (outside `/api/v1`) | No |
| POST | `/api/v1/staff/create` | Create a hospital staff account | No |
| POST | `/api/v1/staff/login` | Log in and get an access token | No |
| POST | `/api/v1/patient/search` | Find one patient of the staff member's hospital, falling back to the hospital's HIS | Yes |

## Authentication

`POST /api/v1/staff/login` returns a JWT. Protected endpoints require:

```
Authorization: Bearer <access_token>
```

The token is an HS256 JWT carrying `staff_id` (also `sub`), `username`, `hospital_id`, `hospital_code`, `iat` and `exp` (lifetime `JWT_TTL`, default `24h`). The **hospital always comes from the token, never from the request body**, so a staff member can only ever reach their own hospital's patients.

On every request the staff member is also looked up in the database by `staff_id` and `hospital_id`, so removing or moving a staff member revokes their token immediately.

Failures:

| Status | Message | When |
|---|---|---|
| `401` | `missing or malformed bearer token` | No `Authorization: Bearer <token>` header |
| `401` | `token expired` | `exp` has passed |
| `401` | `invalid token` | Bad signature, wrong algorithm, or no `exp` |
| `401` | `staff is not a member of this hospital` | The staff member no longer exists at the token's hospital |
| `500` | `internal server error` | The membership lookup failed |

---

## 1. GET /healthz

No auth. Returns `200`:

```json
{ "status": "ok" }
```

---

## 2. POST /api/v1/staff/create

Creates a staff account. No login required.

**Request**

```json
{
  "username": "nurse01",
  "password": "P@ssw0rd!",
  "hospital": "hospital-a"
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `username` | string | Yes | Unique within the hospital; the same username may exist at another hospital |
| `password` | string | Yes | Minimum 8 characters, stored as a bcrypt hash |
| `hospital` | string | Yes | Hospital `code`, e.g. `hospital-a` |

**Response `201 Created`**

```json
{
  "id": "0f9c2f1e-...",
  "username": "nurse01",
  "hospital_id": "3b8e...",
  "created_at": "2026-09-18T10:00:00Z"
}
```

**Errors**

| Status | When |
|---|---|
| `400` | The body isn't valid JSON, a field is missing, or the password is shorter than 8 characters |
| `404` | The hospital code doesn't exist |
| `409` | The username already exists at that hospital |
| `500` | Unexpected error |

---

## 3. POST /api/v1/staff/login

Authenticates a staff member. No login required.

**Request**

```json
{
  "username": "nurse01",
  "password": "P@ssw0rd!",
  "hospital": "hospital-a"
}
```

**Response `200 OK`**

```json
{ "access_token": "eyJhbGciOiJIUzI1NiIs..." }
```

**Errors**

| Status | When |
|---|---|
| `400` | The body isn't valid JSON, or a field is missing |
| `401` | Wrong username, password or hospital |
| `429` | More than 5 requests per second from the same IP (burst 10); returned by nginx, not in the JSON error format |
| `500` | Unexpected error |

A wrong username, a wrong hospital and a wrong password all return the same `401`, so the response doesn't reveal which accounts exist.

---

## 4. POST /api/v1/patient/search

Finds the **one** patient registered at the logged-in staff member's hospital that matches the filters. **Login required.**

It's a POST rather than a GET so that national IDs, passport numbers and names stay out of URLs and access logs.

**Request** — every field is optional, but at least one is required. Values are trimmed of surrounding whitespace, and a field that is empty after trimming counts as not given. An empty body, `{}`, or a body whose fields are all blank is rejected with `400`.

```json
{
  "national_id": "1234567890123",
  "passport_id": "AA1234567",
  "first_name": "Somchai",
  "middle_name": "",
  "last_name": "Jaidee",
  "date_of_birth": "1990-05-17",
  "phone_number": "0812345678",
  "email": "somchai@example.com"
}
```

| Field | Type | Notes |
|---|---|---|
| `national_id`, `passport_id` | string | Exact match |
| `first_name`, `middle_name`, `last_name` | string | Matches the Thai **or** the English column |
| `date_of_birth` | string | `YYYY-MM-DD` |
| `phone_number`, `email` | string | Exact match, email is case-insensitive |

Filters are combined with AND. The hospital filter is always applied on top, from the token. When the filters match more than one patient, the response is `409` and the client should add filters (e.g. `national_id`). The HIS is not asked in that case.

**Response `200 OK`** — a single patient object (not an array). `404` when nothing matches.

```json
{
  "first_name_th": "สมชาย",
  "middle_name_th": "",
  "last_name_th": "ใจดี",
  "first_name_en": "Somchai",
  "middle_name_en": "",
  "last_name_en": "Jaidee",
  "date_of_birth": "1990-05-17",
  "patient_hn": "HN0001",
  "national_id": "1234567890123",
  "passport_id": "",
  "phone_number": "0812345678",
  "email": "somchai@example.com",
  "gender": "M"
}
```

`patient_hn`, `phone_number` and `email` come from the `hospital_patients` row, so each hospital sees its own contact details. Everything else comes from `patients`.

**HIS fallback**

A patient found locally is returned as stored; the HIS is **not** called. Only when no local patient matches, the request carries `national_id` or `passport_id`, and the staff member's hospital has an HIS integration (hospital A: `GET {HIS_HOSPITAL_A_SEARCH_PATIENT_URL}/{id}`, e.g. `https://hospital-a.api.co.th/patient/search/{id}`), the middleware asks that HIS:

1. It searches by `national_id` first, then by `passport_id` if the HIS answers `404` for the national ID. A record that doesn't carry the ID that was asked for is rejected (`502`).
2. It stores the record. The HIS is the source of truth: when another hospital already registered the same national ID or passport ID, that shared `patients` row is **replaced** by the HIS record, including values the HIS leaves empty. Only the national ID and passport ID pointing to two *different* stored patients is refused (`409`), since that would merge two people. A `hospital_patients` row is added, or its HN, phone and email are refreshed, for this hospital. Both writes happen in one transaction.
3. It returns the record only if it matches **every** filter in the request; otherwise the response is `404`. The record stays stored either way.

A staff member only ever triggers their own hospital's HIS. A search without an ID, an HIS that doesn't know the patient, or a hospital without an HIS integration all return `404`.

**Errors**

| Status | When |
|---|---|
| `400` | The body isn't valid JSON, has no filter (`at least one search filter is required`), or `date_of_birth` is not `YYYY-MM-DD` |
| `401` | See [Authentication](#authentication) |
| `404` | `patient not found`: no match locally, and the HIS was not asked, doesn't know the patient, or returned a record that fails a filter |
| `409` | `more than one patient matches, add more filters`: the filters match several stored patients |
| `409` | `patient record conflicts with an existing patient`: the HIS record's national ID and passport ID belong to different stored patients, or its HN is taken by another patient at this hospital |
| `500` | Unexpected error |
| `502` | The HIS is unreachable, answered with a status other than `200`/`404`, or returned a body that can't be decoded, can't be stored (no ID, bad date of birth, unknown gender, no HN) or is for a different ID |
| `504` | The HIS timed out (`HIS_HOSPITAL_A_TIMEOUT`, default `10s`) |

---

## Error format

Every error uses the same envelope (`handler.ErrorResponse`):

```json
{
  "code": 400,
  "message": "invalid date_of_birth \"17-05-1990\": expected YYYY-MM-DD"
}
```

Messages describe what the client should fix. They never include a password, a token or a database error string; unexpected errors return `500` with `internal server error`. The only response outside this format is nginx's `429` on login.
