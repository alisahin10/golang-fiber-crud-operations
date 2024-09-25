# Go Fiber User Management API

This project is a User Management API built with [Go Fiber](https://gofiber.io/) framework. It provides RESTful endpoints for creating and managing users, utilizing BuntDB for storage and Zap for logging.

## Features

- **User Creation**: Register new users with validation and password hashing.
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes.
- **Logging**: Structured logging using Zap, with options to log to a file in JSON format.
- **Middleware**: Custom middleware for request logging and processing.
- **Environment Configuration**: Load environment variables from a `.env` file.

## Technologies Used

- Go
- Fiber
- BuntDB
- Zap Logger
- Godotenv
- bcrypt

## Project Structure
![image](https://github.com/user-attachments/assets/3acecf22-b9f9-46e0-915a-9886e0bb06c6)



