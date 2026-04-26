ALTER TABLE links ADD COLUMN folder_id INTEGER REFERENCES folders(id) ON DELETE SET NULL;
CREATE INDEX idx_links_folder_id ON links(folder_id);
