CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    shortened_url VARCHAR(10) NOT NULL UNIQUE
);
