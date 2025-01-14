CREATE TABLE comments
(
    id SERIAL PRIMARY KEY,
    news_id INT DEFAULT NULL, 
    parent_comment_id INT DEFAULT NULL, 
    content TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE comments OWNER TO admin;
