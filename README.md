User Management API with Go Fiber and BuntDB
This project is a basic CRUD (Create, Read, Update, Delete) API for user management, built using the Go Fiber web framework and BuntDB, an embedded key/value store. The API provides functionality for creating, reading, updating, deleting, and searching users. Passwords are securely hashed before being stored, and sensitive information like passwords is not exposed in API responses.

Features
Create User: Register a new user with name, email, and password.
Read User: Retrieve user details by user ID.
Update User: Modify existing user details by user ID.
Delete User: Remove a user from the system by user ID.
Get All Users: Retrieve all registered users.
Search Users: Search for users by name or email.
Password Security: User passwords are securely hashed and not included in API responses.
Technologies
Go: The programming language used to build the API.
Fiber: A fast, minimalist web framework inspired by Express.js.
BuntDB: An embeddable, in-memory, fast key/value store with optional persistence.
bcrypt: Used for securely hashing user passwords.
