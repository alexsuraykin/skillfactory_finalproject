CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    pub_date VARCHAR(255) DEFAULT '',
    link TEXT NOT NULL UNIQUE
);

ALTER TABLE feeds OWNER TO admin;