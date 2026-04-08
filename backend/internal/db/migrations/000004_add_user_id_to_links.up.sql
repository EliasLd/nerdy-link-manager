ALTER TABLE links ADD COLUMN user_id INTEGER;

UPDATE links
SET user_id = (
  SELECT id FROM users ORDER BY id ASC LIMIT 1
)
WHERE user_id IS NULL;

CREATE INDEX idx_links_user_id ON links(user_id);
