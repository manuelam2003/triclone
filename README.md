# Group Expense Management API

This API is designed to manage group expenses, settlements, users, and memberships. It supports various operations, including creating users, groups, tracking group expenses, managing settlements between users, and authenticating users.

## Table of Contents
- [Installation](#installation)
- [Endpoints](#endpoints)
- [Authentication](#authentication)
- [Middleware](#middleware)
- [Configuration](#configuration)
- [Example](#example-api-workflow)
- [Table schema](#table-schema)

---

## Installation

1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```

2. Navigate to the directory and install the dependencies:
   ```bash
   cd <project_directory>
   go mod download
   ```

3. Set up the database and ensure the environment variables are set up correctly.
    ```bash
    make db/reset
    ```

4. Run the server:
   ```bash
   make run/api
   ```

---

## Endpoints

### Health Check

- **GET** `/v1/healthcheck`: Checks if the API is running correctly.

### Users

- **GET** `/v1/users`: Retrieve a list of users.

- **GET** `/v1/users/:user_id`: Retrieve a specific user by their ID.

- **POST** `/v1/users`: Create a new user.

- **PUT** `/v1/users/activated`: Activate a user.

- **PATCH** `/v1/users/:user_id`: Update a specific user's details.

- **DELETE** `/v1/users/:user_id`: Delete a specific user.

### Groups

- **GET** `/v1/groups`: Retrieve a list of groups.

- **GET** `/v1/groups/:group_id`: Retrieve a specific group by its ID.

- **POST** `/v1/groups`: Create a new group.

- **PATCH** `/v1/groups/:group_id`: Update group details.

- **DELETE** `/v1/groups/:group_id`: Delete a specific group.

### Group Memberships

- **GET** `/v1/groups/:group_id/members`: Retrieve all members of a group.

- **POST** `/v1/groups/:group_id/members`: Add a new member to a group.

- **DELETE** `/v1/groups/:group_id/members/:user_id`: Remove a member from a group.

- **PUT** `/v1/groups/:group_id/members/:user_id`: Reinstate a member of a group.

### Expenses

- **GET** `/v1/groups/:group_id/expenses`: List all expenses for a group.

- **GET** `/v1/groups/:group_id/expenses/:expense_id`: Retrieve a specific expense.

- **POST** `/v1/groups/:group_id/expenses`: Create a new expense.

- **PUT** `/v1/groups/:group_id/expenses/:expense_id`: Update a specific expense.

- **DELETE** `/v1/groups/:group_id/expenses/:expense_id`: Delete a specific expense.

### Expense Participants

- **GET** `/v1/groups/:group_id/expenses/:expense_id/participants`: List all participants of a specific expense.

- **POST** `/v1/groups/:group_id/expenses/:expense_id/participants`: Add participants to a specific expense.

- **PUT** `/v1/groups/:group_id/expenses/:expense_id/participants/:participant_id`: Update a participant in an expense.

- **DELETE** `/v1/groups/:group_id/expenses/:expense_id/participants/:participant_id`: Delete a participant from an expense.

### Settlements

- **GET** `/v1/groups/:group_id/settlements`: List all settlements for a group.

- **GET** `/v1/groups/:group_id/settlements/:settlement_id`: Retrieve a specific settlement.

- **POST** `/v1/groups/:group_id/settlements`: Add a new settlement.

- **DELETE** `/v1/groups/:group_id/settlements/:settlement_id`: Delete a specific settlement.

### Authentication

- **POST** `/v1/tokens/authentication`: Authenticate a user and create an authentication token.

---

## Authentication

Authentication is handled using JWT tokens. To access secured routes, users must include a valid JWT token in the `Authorization` header in the format:

```
Authorization: Bearer <token>
```

Tokens are generated via the `/v1/tokens/authentication` endpoint after a user successfully logs in.

## Middleware

- **Panic Recovery**: Catches and recovers from any unexpected server panics.
- **Rate Limiting**: Controls the rate at which requests can be made to the API.
- **Authentication**: Ensures only authenticated users can access certain routes.
- **Require Activation**: Some routes require that the user is activated before they can access them.

## Configuration

- The API requires a `.env` file or configuration management for settings like database connections and JWT secret keys.
- Ensure that the environment variables are set for running the server in production.

## Example API Workflow

1. **User Registration/Login**: Use `/auth/signup` and `/auth/login` for user registration and login.
2. **Group Management**: Create a group using `/groups`, and add members with `/groups/{group_id}/members`.
3. **Expense Management**: Add expenses using `/groups/{group_id}/expenses` and assign participants with `/expenses/{expense_id}/participants`.
4. **Settlements**: When debts are settled, use `/groups/{group_id}/settlements`.
5. **Balance Check**: Use `/groups/{group_id}/balance` to see who owes whom and how much.

## Table schema

### 1. **Users Table**

- **user_id** (Primary Key): Unique identifier for each user.
- **name**: The user's name.
- **email**: The user's email address (unique).
- **password**: Hashed password for authentication.
- **created_at**: Timestamp when the user was created.

### 2. **Groups Table**

- **group_id** (Primary Key): Unique identifier for each group.
- **group_name**: Name of the group.
- **created_by** (Foreign Key -> Users): The user who created the group.
- **created_at**: Timestamp when the group was created.

### 3. **Group Members Table**

- **group_member_id** (Primary Key): Unique identifier for each group member record.
- **group_id** (Foreign Key -> Groups): The group to which the user belongs.
- **user_id** (Foreign Key -> Users): The user who is a member of the group.
- **joined_at**: Timestamp when the user joined the group.
- **is_active**: Indicates if the user is still part of the group.
- **left_at**: Timestamp when the user left the group.

### 4. **Expenses Table**

- **expense_id** (Primary Key): Unique identifier for each expense.
- **group_id** (Foreign Key -> Groups): The group associated with the expense.
- **amount**: The total amount of the expense.
- **description**: Description of the expense (e.g., "Dinner").
- **paid_by** (Foreign Key -> Users): The user who paid for the expense.
- **created_at**: Timestamp when the expense was created.

### 5. **Expense Participants Table**

- **expense_participant_id** (Primary Key): Unique identifier for each participant record.
- **expense_id** (Foreign Key -> Expenses): The expense associated with this record.
- **user_id** (Foreign Key -> Users): The user who is participating in the expense.
- **amount_owed**: The amount that this user owes for the expense.

### 6. **Settlements Table**

- **settlement_id** (Primary Key): Unique identifier for each settlement.
- **group_id** (Foreign Key -> Groups): The group associated with the settlement.
- **payer_id** (Foreign Key -> Users): The user who is making the payment to settle debts.
- **payee_id** (Foreign Key -> Users): The user who is receiving the payment.
- **amount**: The amount being settled.
- **settled_at**: Timestamp when the settlement occurred.# Group Expense Management API
