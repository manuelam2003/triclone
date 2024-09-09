-- Down Migration

-- 1. Delete from settlements first since it references users and groups
DELETE FROM settlements WHERE group_id IN (1, 2);

-- 2. Delete from expense_participants as it references users and expenses
DELETE FROM expense_participants WHERE expense_id IN (1, 2, 3);

-- 3. Delete from expenses since it references groups and users
DELETE FROM expenses WHERE group_id IN (1, 2);

-- 4. Delete from group_members as it references groups and users
DELETE FROM group_members WHERE group_id IN (1, 2);

-- 5. Delete from groups since it references users
DELETE FROM groups WHERE group_id IN (1, 2);

-- 6. Finally, delete from users
DELETE FROM users WHERE user_id IN (1, 2, 3, 4);
