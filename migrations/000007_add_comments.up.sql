CREATE TABLE IF NOT EXISTS comments (
                                    id bigserial PRIMARY KEY,
                                    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                    news_id bigint NOT NULL REFERENCES news ON DELETE CASCADE,
                                    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                    content text NOT NULL,
                                    version uuid NOT NULL DEFAULT gen_random_uuid()
);