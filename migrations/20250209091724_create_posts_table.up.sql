CREATE TABLE IF NOT EXISTS posts (
    id CHAR(26) PRIMARY KEY,
    body TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    author_id UUID NOT NULL,
    side_id UUID NOT NULL,
    CONSTRAINT fk__post_user FOREIGN KEY(author_id) REFERENCES users(id),
    CONSTRAINT fk__post_side FOREIGN KEY(side_id) REFERENCES sides(id)
)