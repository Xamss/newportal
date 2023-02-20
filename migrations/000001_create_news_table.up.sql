CREATE TABLE IF NOT EXISTS news (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    abstract text NOT NULL,
    tags text[] NOT NULL,
    author text NOT NULL,
    source text NOT NULL,
    time timestamp(0) with time zone NOT NULL,
    version uuid NOT NULL DEFAULT gen_random_uuid()
);