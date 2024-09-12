-- 001_initial_schema.down.sql

-- 1. Drop Indexes
DROP INDEX IF EXISTS idx_group_settlements;
DROP INDEX IF EXISTS idx_group_expenses;
DROP INDEX IF EXISTS idx_expense_id;
DROP INDEX IF EXISTS idx_group_id;

-- 2. Drop Tables
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS settlements;
DROP TABLE IF EXISTS expense_participants;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;