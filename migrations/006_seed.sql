-- Migration: 006_seed_initial_data.sql
-- Вставка начальных данных в правильном порядке с проверками

-- 1. Вставляем пользователей (только если не существуют)
DO $$ 
BEGIN
    -- Проверяем и вставляем пользователей только если их нет
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'monitor') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('monitor', 'monitor@testsystem.ru', 'monitor', 1, '$2y$12$Dc7wN3TQlym69XcfYtsnkOXmH6wY0RWfLSDnpsZfMlEEkrT1OFSHW', 'Монитор User', 'Test System', 'Монитор', '+7-999-000-00-01');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'moderator') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('moderator', 'moderator@testsystem.ru', 'moderator', 1, '$2y$12$Vqmodk5UMpRqjG0HMbOi4e54R5UffACnh7gMU6obZHBO31uwOv59S', 'Модератор User', 'Test System', 'Модератор', '+7-999-000-00-02');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('admin', 'admin@testsystem.ru', 'admin', 1, '$2y$12$1b6PV2G0iUgrrjw9S642QOJxoHamlLr3hN4ww90co/OSUlwmiUcuu', 'Администратор', 'Test System', 'Administrator', '+7-999-000-00-03');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'agent') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('agent', 'agent@testsystem.ru', 'agent', 1, '$2y$12$Zg.mLS/GaVrrPS84kGHU2uTHlYEul18Iip53w/HHU0.DnFAVGk.TC', 'Агент', 'Test System', 'Агент', '+7-999-000-00-04');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'support') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('support', 'support@testsystem.ru', 'support', 1, '$2y$12$JJ3Ygj72LIJ3iEEBxLnF7ubFMHK/U1iWYK1RvGH7g7EY2cMRDqB9K', 'Техподдержка', 'Test System', 'Техническая поддержка', '+7-999-000-00-05');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'partner') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('partner', 'partner@testsystem.ru', 'partner', 1, '$2y$12$NRo08V/jZNqHjf0w6JwvcOZ0SUrVAdvuA9Lyq.lCC6zchcXxahH2a', 'Партнер', 'Taxi Company', 'Управляющий', '+7-999-000-00-06');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'operator1') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('operator1', 'operator1@testsystem.ru', 'operator', 1, '$2y$12$testpasswordhashforoperator1', 'Оператор 1', 'Test System', 'Оператор', '+7-999-000-00-07');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'tech1') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('tech1', 'tech1@testsystem.ru', 'technician', 1, '$2y$12$testpasswordhashfortechnician1', 'Техник 1', 'Test System', 'Техник', '+7-999-000-00-08');
    END IF;
END $$;

-- 2. Вставляем локации (только если не существуют)
INSERT INTO locations (name, address, contact_person, contact_phone, monthly_rent, rent_due_day) 
SELECT 'ТЦ "Москва"', 'ул. Ленина, 1', 'Иванов Иван', '+7-999-123-45-67', 15000.00, 15
WHERE NOT EXISTS (SELECT 1 FROM locations WHERE name = 'ТЦ "Москва"');

INSERT INTO locations (name, address, contact_person, contact_phone, monthly_rent, rent_due_day) 
SELECT 'ТРК "Европа"', 'пр. Мира, 25', 'Петрова Мария', '+7-999-765-43-21', 20000.00, 10
WHERE NOT EXISTS (SELECT 1 FROM locations WHERE name = 'ТРК "Европа"');

INSERT INTO locations (name, address, contact_person, contact_phone, monthly_rent, rent_due_day) 
SELECT 'Аэропорт', 'ш. Аэропортовское, 10', 'Сидоров Алексей', '+7-999-555-44-33', 30000.00, 5
WHERE NOT EXISTS (SELECT 1 FROM locations WHERE name = 'Аэропорт');

-- 3. Вставляем вендинговые автоматы (только если не существуют)
INSERT INTO vending_machines (serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date, last_maintenance_date, next_maintenance_date) 
SELECT 'VM001', 'ToyMaster 3000', 'active', id, 45, 100, 1500.00, '2024-01-01', '2024-01-10', '2024-02-10'
FROM locations WHERE name = 'ТЦ "Москва"' 
AND NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = 'VM001')
LIMIT 1;

INSERT INTO vending_machines (serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date, last_maintenance_date, next_maintenance_date) 
SELECT 'VM002', 'ToyMaster 3000', 'active', id, 38, 100, 2300.00, '2024-01-02', '2024-01-11', '2024-02-11'
FROM locations WHERE name = 'Аэропорт' 
AND NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = 'VM002')
LIMIT 1;

INSERT INTO vending_machines (serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date, last_maintenance_date, next_maintenance_date) 
SELECT 'VM003', 'ToyMaster 2000', 'active', id, 22, 80, 1800.00, '2024-01-03', '2024-01-12', '2024-02-12'
FROM locations WHERE name = 'ТРК "Европа"' 
AND NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = 'VM003')
LIMIT 1;

INSERT INTO vending_machines (serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date) 
SELECT 'VM004', 'ToyMaster 3000', 'maintenance', id, 0, 100, 0, '2024-01-04'
FROM locations WHERE name = 'ТРК "Европа"' 
AND NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = 'VM004')
LIMIT 1;

INSERT INTO vending_machines (serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date, last_maintenance_date, next_maintenance_date) 
SELECT 'VM005', 'ToyMaster 2000', 'active', id, 65, 80, 3150.00, '2024-01-05', '2024-01-13', '2024-02-13'
FROM locations WHERE name = 'ТЦ "Москва"' 
AND NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = 'VM005')
LIMIT 1;

-- 4. Вставляем операции (только если не существуют для избежания дубликатов)
INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 1, 'restock', id, '2024-01-15 09:00:00', 10, 60, 50, 1500.00, 1500.00, 0, 'Пополнение игрушек'
FROM users WHERE username = 'admin' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 1 AND operation_type = 'restock' AND operation_date = '2024-01-15 09:00:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 2, 'restock', id, '2024-01-15 10:30:00', 15, 65, 50, 2300.00, 2300.00, 0, 'Пополнение игрушек'
FROM users WHERE username = 'admin' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 2 AND operation_type = 'restock' AND operation_date = '2024-01-15 10:30:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 3, 'restock', id, '2024-01-15 11:15:00', 8, 58, 50, 1800.00, 1800.00, 0, 'Пополнение игрушек'
FROM users WHERE username = 'operator1' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 3 AND operation_type = 'restock' AND operation_date = '2024-01-15 11:15:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 1, 'collection', id, '2024-01-14 16:00:00', 25, 25, 0, 3200.00, 200.00, 3000.00, 'Инкассация денежных средств'
FROM users WHERE username = 'admin' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 1 AND operation_type = 'collection' AND operation_date = '2024-01-14 16:00:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 2, 'collection', id, '2024-01-14 16:30:00', 30, 30, 0, 2850.00, 350.00, 2500.00, 'Инкассация денежных средств'
FROM users WHERE username = 'admin' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 2 AND operation_type = 'collection' AND operation_date = '2024-01-14 16:30:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 3, 'collection', id, '2024-01-14 17:00:00', 22, 22, 0, 1950.00, 150.00, 1800.00, 'Инкассация денежных средств'
FROM users WHERE username = 'operator1' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 3 AND operation_type = 'collection' AND operation_date = '2024-01-14 17:00:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 1, 'maintenance', id, '2024-01-10 14:00:00', 45, 45, 0, 1500.00, 1500.00, 0, 'Плановое техническое обслуживание'
FROM users WHERE username = 'tech1' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 1 AND operation_type = 'maintenance' AND operation_date = '2024-01-10 14:00:00')
LIMIT 1;

INSERT INTO vending_operations (vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected, notes)
SELECT 2, 'maintenance', id, '2024-01-11 15:30:00', 38, 38, 0, 2300.00, 2300.00, 0, 'Плановое техническое обслуживание'
FROM users WHERE username = 'tech1' 
AND NOT EXISTS (SELECT 1 FROM vending_operations WHERE vending_machine_id = 2 AND operation_type = 'maintenance' AND operation_date = '2024-01-11 15:30:00')
LIMIT 1;