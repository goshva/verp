-- Заполнение тестовыми данными для вендинговых автоматов
-- Вставляем только если локации существуют
INSERT INTO vending_machines 
(serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date)
SELECT 
    'VM001', 'ToyMaster 3000', 'active', id, 45, 100, 1500.00, '2024-01-01'
FROM locations WHERE name = 'ТЦ "Москва"'
LIMIT 1;

INSERT INTO vending_machines 
(serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date)
SELECT 
    'VM002', 'ToyMaster 3000', 'active', id, 38, 100, 2300.00, '2024-01-02'
FROM locations WHERE name = 'БЦ "Сити"'
LIMIT 1;

INSERT INTO vending_machines 
(serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date)
SELECT 
    'VM003', 'ToyMaster 2000', 'active', id, 22, 80, 1800.00, '2024-01-03'
FROM locations WHERE name = 'ТРЦ "Европа"'
LIMIT 1;

INSERT INTO vending_machines 
(serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date)
SELECT 
    'VM004', 'ToyMaster 3000', 'maintenance', id, 0, 100, 0, '2024-01-04'
FROM locations WHERE name = 'Аэропорт "Северный"'
LIMIT 1;

INSERT INTO vending_machines 
(serial_number, model, status, location_id, current_toys_count, capacity_toys, cash_amount, installation_date)
SELECT 
    'VM005', 'ToyMaster 2000', 'active', id, 65, 80, 3150.00, '2024-01-05'
FROM locations WHERE name = 'ЖД Вокзал "Центральный"'
LIMIT 1;