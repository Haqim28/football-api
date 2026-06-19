# Football API — Technical Test AYO 2026

REST API backend untuk manajemen tim sepakbola amatir perusahaan XYZ.

**Stack:** Go 1.21 · Gin · GORM · PostgreSQL · Docker

---

## Swagger UI

Setelah server berjalan, akses:
```
http://localhost:8080/swagger/index.html
```

Klik tombol **Authorize** (kanan atas) → masukkan API Key untuk test endpoint yang butuh auth.

---

## Cara Menjalankan

### Prerequisites
- Docker & Docker Compose, **atau**
- Go 1.21+ dan PostgreSQL lokal

### Dengan Docker (recommended)

```bash
git clone <repo-url>
cd football-api
docker-compose up -d
```

API akan berjalan di `http://localhost:8080`

### Tanpa Docker

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi DB lokal

go mod tidy

# (Opsional) Regenerate swagger docs jika ada perubahan kode:
# go install github.com/swaggo/swag/cmd/swag@latest
# swag init -g cmd/main.go
go run ./cmd/main.go
```

---

## Autentikasi

Endpoint **write** (POST/PUT/DELETE) butuh header:

```
X-API-Key: my-super-secret-api-key
```

Endpoint read (GET) tidak butuh autentikasi.

---

## Endpoints

### Health Check
```
GET /health
```

---

### 1. Tim (Teams)

#### Tambah Tim
```
POST /api/v1/teams
Header: X-API-Key: <key>

Body:
{
  "name": "Persib Bandung",
  "logo": "https://example.com/logo.png",
  "founded": 1933,
  "headquarters_address": "Jl. Sulanjana No.17",
  "headquarters_city": "Bandung"
}
```

#### Daftar Semua Tim
```
GET /api/v1/teams
```

#### Detail Tim
```
GET /api/v1/teams/:id
```

#### Update Tim
```
PUT /api/v1/teams/:id
Header: X-API-Key: <key>

Body: (semua field optional)
{
  "name": "Persib Bandung FC",
  "headquarters_city": "Kota Bandung"
}
```

#### Hapus Tim (Soft Delete)
```
DELETE /api/v1/teams/:id
Header: X-API-Key: <key>
```

---

### 2. Pemain (Players)

#### Tambah Pemain
```
POST /api/v1/players
Header: X-API-Key: <key>

Body:
{
  "team_id": 1,
  "name": "Achmad Jufriyanto",
  "height": 178,
  "weight": 72,
  "position": "bertahan",
  "jersey_number": 5
}
```
> **Posisi valid:** `penyerang` | `gelandang` | `bertahan` | `penjaga_gawang`
> **Nomor punggung** harus unik dalam satu tim (1–99)

#### Daftar Pemain
```
GET /api/v1/players
GET /api/v1/players?team_id=1   # filter per tim
```

#### Detail Pemain
```
GET /api/v1/players/:id
```

#### Update Pemain
```
PUT /api/v1/players/:id
Header: X-API-Key: <key>
```

#### Hapus Pemain (Soft Delete)
```
DELETE /api/v1/players/:id
Header: X-API-Key: <key>
```

---

### 3. Jadwal Pertandingan (Matches)

#### Buat Jadwal
```
POST /api/v1/matches
Header: X-API-Key: <key>

Body:
{
  "match_date": "2026-07-15T15:00:00Z",
  "home_team_id": 1,
  "away_team_id": 2
}
```

#### Daftar Pertandingan
```
GET /api/v1/matches
```

#### Detail Pertandingan
```
GET /api/v1/matches/:id
```

#### Update Jadwal
```
PUT /api/v1/matches/:id
Header: X-API-Key: <key>

Body:
{
  "match_date": "2026-07-20T19:00:00Z"
}
```
> Tidak bisa update pertandingan yang sudah `completed`

#### Hapus Jadwal (Soft Delete)
```
DELETE /api/v1/matches/:id
Header: X-API-Key: <key>
```

---

### 4. Hasil Pertandingan

```
POST /api/v1/matches/:id/result
Header: X-API-Key: <key>

Body:
{
  "home_score": 2,
  "away_score": 1,
  "goals": [
    { "player_id": 3, "minute": 23 },
    { "player_id": 3, "minute": 67 },
    { "player_id": 8, "minute": 45 }
  ]
}
```

> **Validasi:**
> - Pemain yang cetak gol harus dari salah satu tim yang bermain
> - Jumlah gol dari player home_team harus = `home_score`
> - Jumlah gol dari player away_team harus = `away_score`
> - Setelah submit, status match otomatis jadi `completed`

---

### 5. Report

#### Semua Laporan
```
GET /api/v1/reports
```

#### Laporan per Pertandingan
```
GET /api/v1/reports/matches/:id
```

**Contoh Response:**
```json
{
  "success": true,
  "message": "laporan berhasil diambil",
  "data": {
    "match_id": 1,
    "match_date": "2026-07-15T15:00:00Z",
    "home_team": { "id": 1, "name": "Persib Bandung", "logo": "..." },
    "away_team": { "id": 2, "name": "Persija Jakarta", "logo": "..." },
    "home_score": 2,
    "away_score": 1,
    "status": "completed",
    "match_result": "Tim Home Menang",
    "top_scorer": {
      "player_id": 3,
      "name": "Achmad Jufriyanto",
      "team_name": "Persib Bandung",
      "goals": 2
    },
    "home_team_total_wins": 3,
    "away_team_total_wins": 1,
    "goals": [
      { "player_id": 3, "player_name": "Achmad Jufriyanto", "team_name": "Persib Bandung", "minute": 23 },
      { "player_id": 8, "player_name": "Marko Simic", "team_name": "Persija Jakarta", "minute": 45 },
      { "player_id": 3, "player_name": "Achmad Jufriyanto", "team_name": "Persib Bandung", "minute": 67 }
    ]
  }
}
```

---

## Format Response

Semua response menggunakan format konsisten:

```json
{
  "success": true,
  "message": "...",
  "data": { ... }
}
```

Error:
```json
{
  "success": false,
  "message": "request failed",
  "error": "deskripsi error"
}
```

---

## Asumsi

1. **Autentikasi** menggunakan API Key via header `X-API-Key` — melindungi semua operasi write (POST/PUT/DELETE). Read (GET) bersifat publik.
2. **Soft delete** diimplementasikan via GORM `DeletedAt` (field `deleted_at`) untuk semua entitas.
3. **Nomor punggung** unik per tim, bukan global.
4. **Submit hasil pertandingan** bersifat idempotent — jika match belum `completed`, bisa di-submit ulang (gol lama dihapus, gol baru diinsert). Match yang sudah `completed` tidak bisa diubah.
5. **Akumulasi kemenangan** dihitung dari seluruh pertandingan yang sudah selesai dengan ID ≤ match yang diminta (asumsi ID berurutan = urutan pertandingan dimainkan).
6. **Top scorer** adalah pemain dengan gol terbanyak di pertandingan tersebut. Jika ada seri, yang pertama ditemukan yang diambil.
