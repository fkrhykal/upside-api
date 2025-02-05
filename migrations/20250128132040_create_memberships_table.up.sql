CREATE TYPE role AS ENUM ('founder', 'admin', 'member');
CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY,
    member_id UUID NOT NULL,
    side_id UUID NOT NULL,
    role role NOT NULL,
    CONSTRAINT fk__membership_user FOREIGN KEY(member_id) REFERENCES users(id),
    CONSTRAINT fk__membership_side FOREIGN KEY(side_id) REFERENCES sides(id)
);