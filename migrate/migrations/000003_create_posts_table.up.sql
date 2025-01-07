CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER REFERENCES users(id),
                       content TEXT NOT NULL
);

