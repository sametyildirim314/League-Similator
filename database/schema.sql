-- Teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER REFERENCES teams(id),
    away_team_id INTEGER REFERENCES teams(id),
    home_score INTEGER,
    away_score INTEGER,
    week INTEGER NOT NULL,
    played BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- League table
CREATE TABLE IF NOT EXISTS league_table (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) UNIQUE,
    points INTEGER DEFAULT 0,
    played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    goal_difference INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Predictions table
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    predicted_position INTEGER NOT NULL,
    predicted_points INTEGER NOT NULL,
    prediction_percentage DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default teams if they don't exist
INSERT INTO teams (name)
SELECT 'Manchester United' WHERE NOT EXISTS (SELECT 1 FROM teams WHERE name = 'Manchester United');

INSERT INTO teams (name)
SELECT 'Liverpool' WHERE NOT EXISTS (SELECT 1 FROM teams WHERE name = 'Liverpool');

INSERT INTO teams (name)
SELECT 'Chelsea' WHERE NOT EXISTS (SELECT 1 FROM teams WHERE name = 'Chelsea');

INSERT INTO teams (name)
SELECT 'Arsenal' WHERE NOT EXISTS (SELECT 1 FROM teams WHERE name = 'Arsenal');

-- Create initial league table entries for each team if they don't exist
INSERT INTO league_table (team_id)
SELECT id FROM teams WHERE NOT EXISTS (SELECT 1 FROM league_table WHERE team_id = teams.id); 