CREATE TABLE promotion (
    id SERIAL PRIMARY KEY,
    amount VARCHAR(255) DEFAULT '50',
    days INT DEFAULT 7
);
