CREATE INDEX IF NOT EXISTS news_title_idx ON news USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS news_author_idx ON news USING GIN (to_tsvector('simple', author));
CREATE INDEX IF NOT EXISTS news_tags_idx ON news USING GIN (tags);