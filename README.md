# Premier League Simulation

A REST API for simulating a Premier League season with four teams, following Premier League rules.

## Features

- League structure with four teams
- Match simulation and league table generation
- Championship probability prediction
- Complete API for managing the simulation
- System reset functionality to restart the simulation
- Automatic fixture generation
- Sequential week simulation (previous weeks must be simulated first)
- Automatic championship predictions after week 4 and on all subsequent week simulations

## Requirements

- Go 1.19+
- PostgreSQL 14+
- Docker and Docker Compose (recommended for easy setup)

## Setup

### Using Docker (Recommended)

The easiest way to run the application is using Docker Compose:

```bash
# Start all services (PostgreSQL, pgAdmin and the application)
docker compose up -d

# Stop and remove services
docker compose down
```

The API will be available at `http://localhost:8081`.


## Database Access

- pgAdmin will be available at `http://localhost:5050`:
  - Email: admin@admin.com
  - Password: admin

To connect to the database through pgAdmin:

1. Login to pgAdmin
2. Right-click on "Servers" and select "Register > Server"
3. On the "General" tab, give it a name (e.g., "Premier League")
4. On the "Connection" tab, enter:
   - Host: postgres (if connecting from Docker) or localhost (if connecting from host)
   - Port: 5432
   - Maintenance database: premier_league
   - Username: postgres
   - Password: postgres
5. Click "Save"

## API Endpoints

### Teams

- `GET /api/teams` - List all teams
- `GET /api/teams/:id` - Get team details

### Matches

- `GET /api/matches` - List all matches
- `GET /api/matches/week/:week` - Get matches for a specific week
- `POST /api/matches/simulate/:week` - Simulate matches for a specific week (automatically generates fixtures if needed)
  - Note: You must simulate weeks in order (week 1, then week 2, etc.)
- `POST /api/matches/simulate-all` - Simulate all remaining matches (automatically generates fixtures if needed)

### League

- `GET /api/league/table` - Get current league table
- `GET /api/league/table/week/:week` - Get league table at a specific week

### Predictions

- `GET /api/predictions` - Get current championship predictions
  - Note: Predictions are automatically generated after simulating week 4 and updated after each subsequent week simulation

### System

- `POST /api/system/reset` - Reset the entire system (clear matches, reset league table, delete predictions)


All API endpoints can be easily tested using Postman or any other API client:

1. Download and install [Postman](https://www.postman.com/downloads/)
2. Create a new request with the appropriate HTTP method (GET, POST)
3. Enter the URL for the endpoint you want to test (e.g., `http://localhost:8081/api/teams`)
4. For POST requests, no request body is required for this API
5. Click "Send" to execute the request


## System Reset

If you want to reset the system and start a new simulation:

```bash
# Reset the system
curl -X POST http://localhost:8081/api/system/reset

# Start simulation (will automatically generate fixtures)
curl -X POST http://localhost:8081/api/matches/simulate-all
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
