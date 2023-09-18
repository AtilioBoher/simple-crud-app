CREATE TABLE IF NOT EXIST users (
    user_id INT AUTO_INCREMENT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(50),
    user_age INT,
    PRIMARY KEY (user_id)
);

INSERT INTO users (user_name, user_email, user_age) values ('Fulano', 'fulano@mail.com', 18);
INSERT INTO users (user_name, user_email, user_age) values ('Mengano', 'mengano@mail.com', 20);
