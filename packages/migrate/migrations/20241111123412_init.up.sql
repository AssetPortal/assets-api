-- assets
CREATE TABLE IF NOT EXISTS assets (
    _id SERIAL PRIMARY KEY,                
    id TEXT NOT NULL UNIQUE,         
    address TEXT NOT NULL,               
    description TEXT NULL,                
    image TEXT NULL,                      
    social JSONB NULL,                     
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ                
);

CREATE INDEX idx_asset_id ON assets (id);

-- tokens
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE
);
