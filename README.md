# Pokemon Chatbot API

**Status**: ✅ Complete & Deployed

A Golang-based backend API for the Pokemon Information Chatbot, built for the Kata.ai Software Engineer Enterprise assignment.

## Live Demo

- **API URL**: https://pokemon-api-production-3864.up.railway.app
- **Health Check**: https://pokemon-api-production-3864.up.railway.app/health
- **Telegram Bot**: [@jeko_pokemon_bot](https://t.me/jeko_pokemon_bot)

## Features

- User registration with Telegram ID and name validation
- Pokemon information lookup via PokeAPI (by name or ID)
- PostgreSQL database with Supabase (REST API)
- RESTful API with Gin framework
- Image sprite URLs for Pokemon display
- Production deployment on Railway
- **Dashboard APIs**: Paginated users list, search statistics
- **Search tracking**: Pokemon searches logged to database

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| Framework | Gin (HTTP) |
| Database | Supabase (PostgreSQL via REST API) |
| External API | PokeAPI |
| Deployment | Railway |

## API Endpoints

### Health Check
```
GET /health

Response:
{
  "status": "ok",
  "timestamp": "2025-11-25T15:00:00Z"
}
```

### User Registration
```
POST /api/users/register
Content-Type: application/json

{
  "telegram_id": "123456789",
  "first_name": "John",
  "last_name": "Doe",      // optional
  "username": "johndoe"    // optional
}

Response:
{
  "success": true,
  "message": "User registered successfully",
  "user": { ... }
}
```

### Get User by Telegram ID
```
GET /api/users/:telegramId
```

### Check Registration Status
```
GET /api/users/:telegramId/check

Response:
{
  "registered": true
}
```

### Get Pokemon Information
```
GET /api/pokemon/:name

Examples:
- GET /api/pokemon/pikachu  (by name)
- GET /api/pokemon/25       (by ID)

Response (Found):
{
  "found": true,
  "message": "Pikachu is an <Electric> type Pokemon with 6.0 weight and 0.4 height, here's a picture of Pikachu.",
  "data": {
    "id": 25,
    "name": "Pikachu",
    "types": "Electric",
    "abilities": "Static, Lightning rod",
    "stats": {
      "hp": 35,
      "attack": 55,
      "defense": 40,
      "spAttack": 50,
      "spDefense": 50,
      "speed": 90
    },
    "height": "0.4",
    "weight": "6.0",
    "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/25.png"
  }
}

Response (Not Found):
{
  "found": false,
  "message": "Sorry we don't have information for <fakemon>"
}
```

### List Users (Paginated) - Dashboard API
```
GET /api/users?page=1&limit=10

Response:
{
  "success": true,
  "data": {
    "users": [...],
    "total": 39,
    "page": 1,
    "limit": 10,
    "total_pages": 4
  }
}
```

### Search Statistics - Dashboard API
```
GET /api/stats/searches

Response:
{
  "success": true,
  "stats": {
    "total_searches": 15,
    "found_searches": 12,
    "not_found_searches": 3,
    "top_searched": [{"pokemon_name": "Pikachu", "count": 5}],
    "recent_searches": [...]
  }
}
```

## Project Structure

```
pokemon-chatbot-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration, Supabase client
│   ├── handlers/
│   │   ├── user_handler.go      # User API handlers (register, list, paginate)
│   │   └── pokemon_handler.go   # Pokemon API handlers (search, stats)
│   ├── repository/
│   │   ├── supabase_client.go   # Supabase REST client
│   │   ├── user_repository.go   # User data access
│   │   └── search_repository.go # Search tracking data access
│   ├── services/
│   │   ├── user_service.go      # User business logic
│   │   └── pokemon_service.go   # Pokemon + search logging
│   └── models/
│       └── user.go              # User model
├── go.mod
├── go.sum
└── README.md
```

## Local Development

### Prerequisites
- Go 1.21+
- Supabase project

### Setup

1. Clone repository:
```bash
git clone https://github.com/jaroteko18/pokemon-api.git
cd pokemon-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create `.env` file:
```env
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
PORT=8080
```

4. Run the server:
```bash
go run cmd/server/main.go
```

5. Test:
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/pokemon/pikachu
```

## Database Schema

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_telegram_id ON users(telegram_id);

-- Pokemon search tracking
CREATE TABLE pokemon_searches (
    id SERIAL PRIMARY KEY,
    pokemon_name VARCHAR(255) NOT NULL,
    pokemon_id INTEGER,
    found BOOLEAN DEFAULT true,
    searched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pokemon_name ON pokemon_searches(pokemon_name);
CREATE INDEX idx_searched_at ON pokemon_searches(searched_at DESC);
```

## Kata Platform Integration

### User Registration Action
```yaml
registerUserApi:
  type: api
  options:
    method: POST
    uri: 'https://pokemon-api-production-3864.up.railway.app/api/users/register'
    body:
      telegram_id: user_$(context.userName)
      first_name: $(context.userName)
      username: $(context.userName)
```

### Pokemon Fetch Action
```yaml
fetchPokemon:
  type: api
  options:
    method: GET
    uri: 'https://pokemon-api-production-3864.up.railway.app/api/pokemon/$(context.pokemonName)'
```

### Using Result in Bot
```yaml
# Show Pokemon description
showPokemonInfo:
  type: text
  options:
    text: $(result.message)

# Show Pokemon image (works on Telegram!)
showPokemonImageOK:
  type: template
  options:
    type: image
    items:
      originalContentUrl: $(result.data.sprite)
      previewImageUrl: $(result.data.sprite)
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| SUPABASE_URL | Supabase project URL | Yes |
| SUPABASE_KEY | Supabase anon/service key | Yes |
| PORT | Server port (default: 8080) | No |

## Deployment (Railway)

The API auto-deploys from GitHub `main` branch.

**Railway Config:**
- Build: `go build -o server cmd/server/main.go`
- Start: `./server`

## Architecture

See [ARCHITECTURE_DIAGRAM.md](./ARCHITECTURE_DIAGRAM.md) for complete system architecture.

```
User (Telegram) → Kata Platform → Backend API (Railway) → Supabase + PokeAPI
```

## Author

**Jarot Eko Saputra**
Software Engineer Enterprise / Bot Builder Assignment
Kata.ai - November 2025
