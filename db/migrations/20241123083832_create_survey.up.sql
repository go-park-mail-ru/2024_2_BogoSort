CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'page_type') THEN
        CREATE TYPE page_type AS ENUM ('mainPage', 
        'advertPage', 
        'advertCreatePage', 
        'cartPage', 
        'categoryPage', 
        'advertEditPage', 
        'userPage', 
        'sellerPage', 
        'searchPage');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS question (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    page page_type NOT NULL,
    trigger_value INT,
    lower_description TEXT NOT NULL,
    upper_description TEXT NOT NULL,
    parent_id UUID NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES question(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS answer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    question_id UUID NOT NULL,
    user_id UUID NOT NULL,
    value INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES question(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
