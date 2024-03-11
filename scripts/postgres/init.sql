CREATE TABLE IF NOT EXISTS  users (
  user_id serial PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(300) UNIQUE NOT NULL,
  password BYTEA NOT NULL, 
  name VARCHAR(128),
  surnname VARCHAR(128),
);
