INSERT INTO user (email, password_hash, created_at, updated_at) 
VALUES (?, ?, NOW(), NOW())
