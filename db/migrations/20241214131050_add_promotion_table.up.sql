CREATE TABLE promotion (
    id SERIAL PRIMARY KEY,
    amount VARCHAR(255) NOT NULL,
    days INT NOT NULL DEFAULT 7
);

INSERT INTO promotion (amount, days) VALUES ('50.00', 7);
