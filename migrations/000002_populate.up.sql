INSERT INTO users (name, email, password_hash)
VALUES
('Alice Johnson', 'alice@example.com', '$2a$12$0lb8HkbQ7eGBE2qcTB/DMuia0fnVDjfJ0dpDXM3oSEQj.YWexpF9G'), -- contra es password
('Bob Smith', 'bob@example.com', '$2a$12$0lb8HkbQ7eGBE2qcTB/DMuia0fnVDjfJ0dpDXM3oSEQj.YWexpF9G'),
('Charlie Brown', 'charlie@example.com', '$2a$12$0lb8HkbQ7eGBE2qcTB/DMuia0fnVDjfJ0dpDXM3oSEQj.YWexpF9G'),
('David Lee', 'david@example.com', '$2a$12$0lb8HkbQ7eGBE2qcTB/DMuia0fnVDjfJ0dpDXM3oSEQj.YWexpF9G');

INSERT INTO groups (name, created_by)
VALUES
('Weekend Trip', 1),
('Office Lunch', 2);

INSERT INTO group_members (group_id, user_id)
VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 2),
(2, 4);

INSERT INTO expenses (group_id, amount, description, paid_by)
VALUES
(1, 300.00, 'Hotel Booking', 1),
(1, 100.00, 'Dinner', 2),
(2, 60.00, 'Lunch', 4);

INSERT INTO expense_participants (expense_id, user_id, amount_owed)
VALUES
(1, 1, 100.00),
(1, 2, 100.00),
(1, 3, 100.00),
(2, 1, 33.33),
(2, 2, 33.33),
(2, 3, 33.33),
(3, 2, 30.00),
(3, 4, 30.00);

INSERT INTO settlements (group_id, payer_id, payee_id, amount)
VALUES
(1, 2, 1, 50.00),
(2, 4, 2, 30.00);
