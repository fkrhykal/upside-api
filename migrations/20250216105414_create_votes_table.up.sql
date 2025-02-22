CREATE TABLE IF NOT EXISTS votes (
    id UUID PRIMARY KEY,
    post_id CHAR(26) NOT NULL,
    voter_id UUID NOT NULL,
    kind SMALLINT NOT NULL CONSTRAINT valid_vote_kind CHECK(kind IN (1, -1)),
    CONSTRAINT fk__vote_post FOREIGN KEY(post_id) REFERENCES posts(id),
    CONSTRAINT fk__vote_user FOREIGN KEY(voter_id) REFERENCES users(id)
)