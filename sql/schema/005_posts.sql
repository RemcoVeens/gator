
-- Add a posts table to the database.

-- A post is a single entry from a feed. It should have:

-- id - a unique identifier for the post
-- created_at - the time the record was created
-- updated_at - the time the record was last updated
-- title - the title of the post
-- url - the URL of the post (this should be unique)
-- description - the description of the post
-- published_at - the time the post was published
-- feed_id - the ID of the feed that the post came from


-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL unique,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE posts;
