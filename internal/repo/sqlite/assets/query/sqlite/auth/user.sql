-- Res: user
-- Table: user
-- GetAll
SELECT * FROM user;

-- Res: user
-- Table: user
-- Get
SELECT * FROM user WHERE id = ?;

-- Res: user
-- Table: user
-- GetByUsername
SELECT * FROM user WHERE username = ?;

-- Res: user
-- Table: user
-- Create
INSERT INTO user (id, short_id, name, username, email, created_by, updated_by, created_at, updated_at)
VALUES (:id, :short_id, :name, :username, :email, :created_by, :updated_by, :created_at, :updated_at);

-- Res: user
-- Table: user
-- Update
UPDATE user
SET
    name = :name,
    username = :username,
    email = :email,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Res: user
-- Table: user
-- Delete
DELETE FROM user WHERE id = ?;