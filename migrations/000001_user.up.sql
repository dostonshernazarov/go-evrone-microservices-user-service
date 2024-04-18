CREATE TABLE user_info (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255),
    username VARCHAR(80),
    email VARCHAR(100),
    password VARCHAR(255),
    bio VARCHAR(255),
    website VARCHAR(255),
    role VARCHAR(50),
    refresh_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);