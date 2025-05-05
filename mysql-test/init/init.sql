DROP DATABASE IF EXISTS fake_db;
CREATE DATABASE fake_db;
USE fake_db;

-- Users table
CREATE TABLE Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash CHAR(60) NOT NULL,
    bio TEXT,
    profile_picture BLOB,
    birth_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    signup_ip VARCHAR(45),
    login_attempts INT DEFAULT 0,
    last_login DATETIME,
    account_balance DECIMAL(10, 2) DEFAULT 0.00,
    user_role ENUM('admin', 'editor', 'user', 'guest') DEFAULT 'user',
    settings JSON,
    phone_number VARCHAR(15),
    email_verified BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- `Groups` table (with backticks)
CREATE TABLE `Groups` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    description TEXT
);

-- Memberships (many-to-many between Users and `Groups`)
CREATE TABLE Memberships (
    user_id INT,
    group_id INT,
    joined_at DATETIME,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (group_id) REFERENCES `Groups`(id)  -- Backticks used here as well
);

-- Posts table
CREATE TABLE Posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    content TEXT,
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Comments table
CREATE TABLE Comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    post_id INT,
    user_id INT,
    content TEXT,
    created_at DATETIME,
    FOREIGN KEY (post_id) REFERENCES Posts(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Insert fake Users with varied nullable fields
DELIMITER $$

CREATE PROCEDURE populate_users()
BEGIN
  DECLARE i INT DEFAULT 1;
  WHILE i <= 10000 DO
    INSERT INTO Users (
      name,
      email,
      password_hash,
      bio,
      profile_picture,
      birth_date,
      is_active,
      signup_ip,
      login_attempts,
      last_login,
      account_balance,
      user_role,
      settings,
      phone_number,
      email_verified,
      created_at,
      updated_at
    )
    VALUES (
      CONCAT('User', i),
      CONCAT('user', i, '@example.com'),
      LPAD(HEX(RAND()*100000000000000000), 60, '0'),
      IF(RAND() > 0.2, CONCAT('This is the bio of User', i), NULL), -- 80% chance bio is filled
      NULL, -- Keeping profile_picture always NULL for simplicity
      DATE_SUB(CURDATE(), INTERVAL FLOOR(RAND()*10000) DAY),
      IF(RAND() > 0.1, TRUE, FALSE),
      IF(RAND() > 0.3, INET_NTOA(FLOOR(RAND()*4294967295)), NULL), -- 70% chance to have IP
      FLOOR(RAND()*5),
      IF(RAND() > 0.3, NOW() - INTERVAL FLOOR(RAND()*365) DAY, NULL), -- 70% chance last_login
      ROUND(RAND()*1000, 2),
      ELT(FLOOR(1 + RAND() * 4), 'admin', 'editor', 'user', 'guest'),
      JSON_OBJECT('theme', IF(RAND() > 0.5, 'dark', 'light'), 'language', 'en'),
      IF(RAND() > 0.3, CONCAT('+123456789', LPAD(i, 3, '0')), NULL), -- 70% chance to have phone
      IF(RAND() > 0.2, TRUE, FALSE),
      NOW(),
      NOW()
    );
    SET i = i + 1;
  END WHILE;
END $$

DELIMITER ;

CALL populate_users();

-- Insert fake Groups (with backticks)
DELIMITER $$
CREATE PROCEDURE populate_groups()
BEGIN
  DECLARE i INT DEFAULT 1;
  WHILE i <= 10 DO
    INSERT INTO `Groups` (name, description)
    VALUES (CONCAT('Group', i), CONCAT('Description of Group ', i));
    SET i = i + 1;
  END WHILE;
END $$
DELIMITER ;
CALL populate_groups();

-- Insert fake Memberships
INSERT INTO Memberships (user_id, group_id, joined_at)
SELECT FLOOR(1 + (RAND() * 100)), FLOOR(1 + (RAND() * 10)), NOW()
FROM information_schema.tables LIMIT 100;

-- Insert fake Posts
DELIMITER $$
CREATE PROCEDURE populate_posts()
BEGIN
  DECLARE i INT DEFAULT 1;
  WHILE i <= 100 DO
    INSERT INTO Posts (user_id, content, created_at)
    VALUES (FLOOR(1 + (RAND() * 100)), CONCAT('Post content ', i), NOW());
    SET i = i + 1;
  END WHILE;
END $$
DELIMITER ;
CALL populate_posts();

-- Insert fake Comments
DELIMITER $$
CREATE PROCEDURE populate_comments()
BEGIN
  DECLARE i INT DEFAULT 1;
  WHILE i <= 100 DO
    INSERT INTO Comments (post_id, user_id, content, created_at)
    VALUES (
      FLOOR(1 + (RAND() * 100)),
      FLOOR(1 + (RAND() * 100)),
      CONCAT('Comment content ', i),
      NOW()
    );
    SET i = i + 1;
  END WHILE;
END $$
DELIMITER ;
CALL populate_comments();
