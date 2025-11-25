# Pokemon Chatbot Architecture Diagram

**Status**: Complete & Deployed
**Bot Version**: v1.0.5
**Telegram**: [@jeko_pokemon_bot](https://t.me/jeko_pokemon_bot)
**API**: https://pokemon-api-production-3864.up.railway.app

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACE                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    ┌─────────────┐                                                          │
│    │  Telegram   │  User sends messages via Telegram Bot                    │
│    │    App      │  @jeko_pokemon_bot                                       │
│    └──────┬──────┘                                                          │
│           │                                                                 │
│           ▼                                                                 │
└───────────┼─────────────────────────────────────────────────────────────────┘
            │
            │ HTTPS (Webhook)
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         KATA PLATFORM (Bot Engine)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                         NLU Layer                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐   │   │
│  │  │  Keyword    │  │  Pokemon    │  │   Supermodel NL (NER)       │   │   │
│  │  │   NLU       │  │   NLU       │  │   - Person Name Detection   │   │   │
│  │  └─────────────┘  └─────────────┘  └─────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Conversation Flows                              │   │
│  │                                                                      │   │
│  │  ┌─────────────────────┐      ┌─────────────────────────────────┐    │   │
│  │  │   Greeting Flow     │      │     Pokemon Search Flow         │    │   │
│  │  │                     │      │                                 │    │   │
│  │  │  States:            │      │  States:                        │    │   │
│  │  │  - init             │      │  - init                         │    │   │
│  │  │  - askName          │      │  - askPokemon                   │    │   │
│  │  │  - confirmName      │      │  - fetchInfo                    │    │   │
│  │  │  - registerUser     │      │  - showFound                    │    │   │
│  │  │                     │      │  - showNotFound                 │    │   │
│  │  └─────────────────────┘      └─────────────────────────────────┘    │   │
│  │             │                              │                         │   │
│  └─────────────┼──────────────────────────────┼─────────────────────────┘   │
│                │                              │                             │
│                │ POST /api/users/register     │ GET /api/pokemon/:name      │
│                ▼                              ▼                             │
└────────────────┼──────────────────────────────┼─────────────────────────────┘
                 │                              │
                 │         HTTPS                │
                 ▼                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BACKEND API (Golang + Gin)                               │
│                 Railway: pokemon-api-production-3864                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────┐  ┌───────────────────────────────────┐   │
│  │      User Handler             │  │       Pokemon Handler             │   │
│  │  POST /api/users/register     │  │  GET /api/pokemon/:name           │   │
│  │  GET /api/users/:telegramId   │  │  GET /api/pokemon/search/:query   │   │
│  │  GET /api/users (paginated)   │  │  GET /api/stats/searches          │   │
│  └───────────────┬───────────────┘  └───────────────┬───────────────────┘   │
│                  │                                  │                       │
│                  ▼                                  ▼                       │
│  ┌───────────────────────────────┐  ┌───────────────────────────────────┐   │
│  │       User Service            │  │       Pokemon Service             │   │
│  │  - Register()                 │  │  - GetPokemon()                   │   │
│  │  - GetUserByTelegramID()      │  │  - GetSearchStats()               │   │
│  │  - GetUsersPaginated()        │  │  - LogSearch()                    │   │
│  └───────────────┬───────────────┘  └───────────────┬───────────────────┘   │
│                  │                                  │                       │
└──────────────────┼──────────────────────────────────┼───────────────────────┘
                   │                                  │
                   │ Supabase REST API                │ HTTPS
                   ▼                                  ▼
┌──────────────────────────────────┐  ┌───────────────────────────────────────┐
│      SUPABASE (Database)         │  │         POKEAPI (External)            │
│  PostgreSQL                      │  │   https://pokeapi.co/api/v2           │
│  kwrhzjuufxwedzbsxjjj.supabase.co│  │                                       │
├──────────────────────────────────┤  ├───────────────────────────────────────┤
│                                  │  │                                       │
│  ┌────────────────────────────┐  │  │  - Pokemon data by name/ID            │
│  │      users table           │  │  │  - Types, abilities, stats            │
│  │  - id (SERIAL)             │  │  │  - Height, weight                     │
│  │  - telegram_id (VARCHAR)   │  │  │  - Sprite images                      │
│  │  - username (VARCHAR)      │  │  │                                       │
│  │  - first_name (VARCHAR)    │  │  │                                       │
│  │  - registered_at           │  │  │                                       │
│  │  - last_active             │  │  │                                       │
│  └────────────────────────────┘  │  │                                       │
│                                  │  │                                       │
│  ┌────────────────────────────┐  │  │                                       │
│  │  pokemon_searches table    │  │  │                                       │
│  │  - id (SERIAL)             │  │  │                                       │
│  │  - pokemon_name (VARCHAR)  │  │  │                                       │
│  │  - pokemon_id (INTEGER)    │  │  │                                       │
│  │  - found (BOOLEAN)         │  │  │                                       │
│  │  - searched_at (TIMESTAMP) │  │  │                                       │
│  └────────────────────────────┘  │  │                                       │
│                                  │  │                                       │
└──────────────────────────────────┘  └───────────────────────────────────────┘
```

## Data Flow

### 1. User Registration Flow
```
User → Telegram → Kata Platform → Backend API → Supabase
  │                    │                │            │
  │  "/start"          │                │            │
  │ ───────────────────>                │            │
  │                    │                │            │
  │  "Hello what's your name?"          │            │
  │ <───────────────────                │            │
  │                    │                │            │
  │  "John"            │                │            │
  │ ───────────────────>                │            │
  │                    │                │            │
  │  "So your name is John?"            │            │
  │ <───────────────────                │            │
  │                    │                │            │
  │  "yes"             │                │            │
  │ ───────────────────>                │            │
  │                    │  POST /register│            │
  │                    │ ───────────────>            │
  │                    │                │  INSERT    │
  │                    │                │ ──────────>│
  │                    │                │  Success   │
  │                    │                │ <──────────│
  │                    │  200 OK        │            │
  │                    │ <───────────────            │
  │  "Welcome to PokeBot, John!"        │            │
  │ <───────────────────                │            │
```

### 2. Pokemon Search Flow
```
User → Telegram → Kata Platform → Backend API → PokeAPI
  │                    │                │            │
  │  "pokemon information"              │            │
  │ ───────────────────>                │            │
  │                    │                │            │
  │  "Which Pokemon?"  │                │            │
  │ <───────────────────                │            │
  │                    │                │            │
  │  "pikachu"         │                │            │
  │ ───────────────────>                │            │
  │                    │ GET /pokemon/pikachu        │
  │                    │ ───────────────>            │
  │                    │                │  GET data  │
  │                    │                │ ──────────>│
  │                    │                │  Pokemon   │
  │                    │                │ <──────────│
  │                    │  JSON response │            │
  │                    │ <───────────────            │
  │  "Pikachu is an Electric type..."   │            │
  │  [Pokemon Image Carousel]           │            │
  │ <───────────────────                │            │
  │                    │                │            │
  │  "Which Pokemon?"  │                │            │
  │ <─────────────────── (loop back)    │            │
```

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Bot Engine | Kata Platform | NLU, conversation management, Telegram integration |
| Backend API | Golang + Gin | Business logic, API endpoints |
| Database | Supabase (PostgreSQL) | User data persistence |
| External API | PokeAPI | Pokemon information source |
| Deployment | Railway | Backend API hosting |
| Messaging | Telegram | User interface |

## API Endpoints

### Backend API (Railway)
- `GET /health` - Health check
- `POST /api/users/register` - Register new user
- `GET /api/users/:telegramId` - Get user by Telegram ID
- `GET /api/users/:telegramId/check` - Check registration status
- `GET /api/pokemon/:name` - Get Pokemon by name or ID
- `GET /api/pokemon/search/:query` - Search Pokemon

### Response Format

**Pokemon Found:**
```json
{
  "found": true,
  "message": "Pikachu is an <Electric> type Pokemon with 6.0 weight and 0.4 height...",
  "data": {
    "id": 25,
    "name": "Pikachu",
    "types": "Electric",
    "abilities": "Static, Lightning rod",
    "stats": { "hp": 35, "attack": 55, ... },
    "height": "0.4",
    "weight": "6.0",
    "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/..."
  }
}
```

**Pokemon Not Found:**
```json
{
  "found": false,
  "message": "Sorry we don't have information for <fakemon>"
}
```
