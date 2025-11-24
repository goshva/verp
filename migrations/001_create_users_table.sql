-- Migration: 001_create_users_table.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified_at TIMESTAMP NULL,
    userrole VARCHAR(50) NOT NULL DEFAULT 'user',
    status INTEGER NOT NULL DEFAULT 1,
    lastipaddr VARCHAR(45),
    fullusername VARCHAR(255),
    companyname VARCHAR(255),
    companyrole VARCHAR(255),
    phone VARCHAR(50),
    password VARCHAR(255) NOT NULL,
    remember_token VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Insert default users
INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole) VALUES
('monitor', 'monitor@testsystem.ru', 'monitor', 1, '$2y$12$Dc7wN3TQlym69XcfYtsnkOXmH6wY0RWfLSDnpsZfMlEEkrT1OFSHW', 'eq User', 'Test System', 'Монитор'),
('Модер', 'moderator@testsystem.ru', 'moderator', 1, '$2y$12$Vqmodk5UMpRqjG0HMbOi4e54R5UffACnh7gMU6obZHBO31uwOv59S', 'moderator', 'Test System', 'Модератор'),
('Админ', 'admin@testsystem.ru', 'admin', 1, '$2y$12$1b6PV2G0iUgrrjw9S642QOJxoHamlLr3hN4ww90co/OSUlwmiUcuu', 'Admin User', 'Test System', 'Administrator'),
('Агент', 'agent@testsystem.ru', 'agent', 1, '$2y$12$Zg.mLS/GaVrrPS84kGHU2uTHlYEul18Iip53w/HHU0.DnFAVGk.TC', '', '', ''),
('Поддержка', 'support@testsystem.ru', 'support', 1, '$2y$12$JJ3Ygj72LIJ3iEEBxLnF7ubFMHK/U1iWYK1RvGH7g7EY2cMRDqB9K', 'Помочь и простить', 'Stuff', 'Техническая поддержка'),
('Таксопарк', 'partner@testsystem.ru', 'partner', 1, '$2y$12$NRo08V/jZNqHjf0w6JwvcOZ0SUrVAdvuA9Lyq.lCC6zchcXxahH2a', 'Помочь и простить', 'Taxi Bluze', 'Управляющий');