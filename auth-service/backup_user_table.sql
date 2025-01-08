DROP TABLE IF EXISTS Users;

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,                    
    username VARCHAR(255) NOT NULL UNIQUE,     
    password VARCHAR(255) NOT NULL,             
    role VARCHAR(50) DEFAULT 'USER',         
 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);