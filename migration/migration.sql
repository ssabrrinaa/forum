CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    token VARCHAR(32) NOT NULL UNIQUE,
    expire_time DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    image TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    description TEXT NOT NULL,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    category_id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories_posts_association (
    association_id UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    post_id UUID NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (category_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_user_id_sessions ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_token_sessions ON sessions (token);

CREATE INDEX IF NOT EXISTS idx_user_id_posts ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_title_posts ON posts (title);

CREATE INDEX IF NOT EXISTS idx_post_id_comments ON comments (post_id);
CREATE INDEX IF NOT EXISTS idx_user_id_comments ON comments (user_id);

CREATE INDEX IF NOT EXISTS idx_category_id_association ON categories_posts_association (category_id);
CREATE INDEX IF NOT EXISTS idx_post_id_association ON categories_posts_association (post_id);


-- INSERT INTO categories(category_id, name) VALUES
-- ('a49ed800-3def-416f-8046-54029890abee', 'Lifestyle Hacks'), 
-- ('b8cfe074-6bc4-4445-8702-7de8ff16b05b', 'Feel-Good Stories'),
-- ('21a8f788-c802-4f05-9361-5f075708fea6', 'Quick Recipes'),
-- ('34a940a4-4f8f-4a3f-8d33-03086957e068', 'Product Reviews'),
-- ('096f9a85-bd8e-422e-9905-026466991df9', 'Home Organization'),
-- ('7b415c1c-8f8a-4801-90a8-6006135fe51c', 'DIY Home Decor'),
-- ('8caa27c4-bc90-4ca4-8f62-0cc66242a16d', 'Budget-Friendly Travel'),
-- ('f2e5c475-e643-4e52-a168-40890f46ef46', 'Motivational Quotes'),
-- ('e4cc1f56-08a1-40ac-bbd7-4177f53071b6', 'Funny Anecdotes'),
-- ('13412ac4-e1c0-4fb8-813b-ad55ed73a09d', 'Mindfulness Moments');


