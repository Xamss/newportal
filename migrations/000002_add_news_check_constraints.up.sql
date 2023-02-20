ALTER TABLE news ADD CONSTRAINT news_time_check CHECK (time <= now());
ALTER TABLE news ADD CONSTRAINT tags_length_check CHECK (array_length(tags, 1) BETWEEN 1 AND 5);