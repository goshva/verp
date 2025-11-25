-- Заполнение тестовыми данными для операций
-- Вставляем только если автоматы и пользователи существуют
INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    1, 'restock', id, '2024-01-15 09:00:00', 10, 60, 50, 1500.00, 1500.00, 0
FROM users WHERE username = 'admin' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    2, 'restock', id, '2024-01-15 10:30:00', 15, 65, 50, 2300.00, 2300.00, 0
FROM users WHERE username = 'admin' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    3, 'restock', id, '2024-01-15 11:15:00', 8, 58, 50, 1800.00, 1800.00, 0
FROM users WHERE username = 'operator1' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    1, 'collection', id, '2024-01-14 16:00:00', 25, 25, 0, 3200.00, 200.00, 3000.00
FROM users WHERE username = 'admin' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    2, 'collection', id, '2024-01-14 16:30:00', 30, 30, 0, 2850.00, 350.00, 2500.00
FROM users WHERE username = 'admin' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    3, 'collection', id, '2024-01-14 17:00:00', 22, 22, 0, 1950.00, 150.00, 1800.00
FROM users WHERE username = 'operator1' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    1, 'maintenance', id, '2024-01-10 14:00:00', 45, 45, 0, 1500.00, 1500.00, 0
FROM users WHERE username = 'tech1' LIMIT 1;

INSERT INTO vending_operations 
(vending_machine_id, operation_type, performed_by, operation_date, toys_before, toys_after, toys_added, cash_before, cash_after, cash_collected)
SELECT 
    2, 'maintenance', id, '2024-01-11 15:30:00', 38, 38, 0, 2300.00, 2300.00, 0
FROM users WHERE username = 'tech1' LIMIT 1;