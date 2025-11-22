SELECT user_id, email, password_hash, created_at, updated_at 
FROM user 
WHERE user_id = ?
