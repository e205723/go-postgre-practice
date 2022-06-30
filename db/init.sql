CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    name CHARACTER(255) NOT NULL UNIQUE,
    password CHARACTER(255) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (name, password) VALUES 
('Yoshisaur', 'password'), 
('Kono', 'password2');
