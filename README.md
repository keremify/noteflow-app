# 🚀 Backend-Focused SaaS Application — Go (Golang)

This project is a **backend-driven SaaS application** built with **Go (Golang) using Gin framework**, focusing on **clean architecture, scalable system design, and production-ready authentication patterns**.

The project demonstrates **professional backend development practices**, including **JWT authentication, refresh token flow, role-based authorization (RBAC), and layered architecture design**.

---

# 🏗️ Backend Architecture Overview

The backend follows **Clean Architecture & Layered Architecture principles**, ensuring:

- Separation of concerns  
- High maintainability  
- Testable business logic  
- Scalable system structure  

### High-Level Request Flow

HTTP Request
↓
Router (Gin)
↓
Handler (Controller Layer)
↓
Service (Business Logic Layer)
↓
Repository (Data Access Layer)
↓
Database (MySQL + GORM)


Each layer has a **single responsibility**, providing a clean and predictable data flow.

---

# 📁 Project Structure

/cmd
/internal
/handlers -> HTTP controllers
/services -> Business logic layer
/repositories -> Database access layer
/models -> Domain entities
/middlewares -> JWT & RBAC middleware
/routes -> Central route definitions
/pkg
/utils -> Token & helper utilities
/config -> Environment configuration

---

# 🔁 Request Lifecycle (Router → Handler → Service → Repository)

### Example: Login Flow

POST /auth/login
↓
Router
↓
AuthHandler
↓
AuthService
↓
UserRepository


This architecture:

- Prevents business logic leakage into controllers  
- Keeps handlers thin  
- Centralizes business rules  
- Makes database logic fully isolated  

---

# 🧩 Layer Responsibilities

## Router Layer
- Defines endpoints  
- Applies middleware  
- Manages route grouping  

## Handler Layer
- Handles HTTP requests  
- Performs validation  
- Delegates to service layer  
- Formats HTTP responses  

## Service Layer (Core Logic)
- Authentication & authorization logic  
- Token lifecycle management  
- Refresh token handling  
- Business rules  

## Repository Layer
- Database communication  
- Query abstraction  
- ORM operations  

---

# 🔐 Authentication & Authorization System

### Implemented Features

- JWT Access Token  
- Refresh Token Flow  
- Token expiration handling  
- Role-Based Access Control (RBAC)  

### Token Flow

Login
↓
Access Token (short-lived)
Refresh Token (long-lived)
↓
Access Token expires
↓
Client sends Refresh Token
↓
New Access + Refresh Token issued


This structure improves **security, scalability, and session control**.

---

# 🛡️ Role-Based Access Control (RBAC)

Routes are protected using:

- Authentication middleware  
- Role validation middleware  

This enables:

- Secure admin endpoints  
- Role-based permission management  
- Scalable authorization system  

---

# 🗄️ Database Layer

- MySQL  
- GORM ORM  
- Repository pattern  
- Clean query abstraction  

---

# 🌐 Frontend

A **simple frontend interface** is implemented using **Vanilla JavaScript and basic HTML/CSS**, providing:

- Authentication forms  
- Token handling  
- Basic UI for API testing  

The frontend is intentionally lightweight to keep the **main focus on backend architecture**.

---

# ⚙️ Technology Stack

### Backend
- Go (Golang)
- Gin Web Framework
- JWT Authentication
- MySQL
- GORM ORM
- Clean Architecture
- RESTful API Design

### Frontend
- Vanilla JavaScript
- HTML / CSS

---

# 🎯 Engineering Focus

This project emphasizes:

- Professional backend architecture  
- Secure authentication flows  
- Clean layered design  
- Scalable SaaS backend structure  

---

# 🚀 Future Improvements

- Redis-based token storage  
- Rate limiting middleware  
- Docker & CI/CD pipeline  
- Distributed caching  
- Request tracing & logging  

---

# 🏁 Summary

This project is designed as a **backend-first SaaS architecture showcase**, focusing on **production-level system design** rather than tutorial-based implementations.

