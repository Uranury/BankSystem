# BankSystem - Mock Banking API

> A RESTful API for a mock banking system with user management, account operations, and transaction tracking

## Features

- 🔐 JWT Authentication for secure access
- 💰 Banking operations (Deposit, Withdraw, Transfer)
- 📊 Transaction history tracking
- 👥 User management system with role-based access
- 🗄️ PostgreSQL Database with automated migrations
- 🐳 Fully Dockerized setup
- 🚀 RESTful API endpoints
- ⚡ Built with Go, SQLx, and Gorilla Mux

## Tech Stack

- **Language**: Go 1.23.1
- **Database**: PostgreSQL 16
- **Router**: Gorilla Mux
- **Migration Tool**: golang-migrate
- **Database Library**: SQLx
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Docker and Docker Compose installed on your machine
- Git

## Quick Start

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd BankSystem
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file with your preferred values:
   ```bash
   DATABASE_USER=postgres
   DATABASE_PASSWORD=your_secure_password
   DATABASE_NAME=bankdb
   DATABASE_PORT=5432
   LISTEN_ADDR=:8080
   JWT_SECRET=your_super_secret_jwt_key_here
   SSL_MODE=disable
   ```

3. **Run the application**
   ```bash
   docker-compose up --build
   ```

4. **Access the API**
   - The API will be available at `http://localhost:8080`
   - Database runs on `localhost:5433` (mapped from container port 5432)

## API Endpoints

### Public Endpoints
```
GET    /users           - Get all users
POST   /signup          - Create new user account
POST   /login           - User login (returns JWT token)
```

### Protected Endpoints (Requires JWT Token)
```
GET    /profile         - Get current user's profile information
POST   /withdraw        - Withdraw money from account
POST   /deposit         - Deposit money to account  
POST   /transfer        - Transfer money between accounts
```

### Admin-Only Endpoints (Requires JWT Token + Admin Role)
```
GET    /transactions    - Get all transactions (admin access only)
```

### Authentication & Authorization

For protected endpoints, include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

**User Roles:**
- `user`: Default role for regular users (can access profile, banking operations)
- `admin`: Administrative role (can access all user endpoints + transaction history)

### Example Requests

**Sign Up:**
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "name": "john",
    "email": "john@example.com", 
    "password": "secure_password"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password"
  }'
```

**Get Profile (requires JWT):**
```bash
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Deposit (requires JWT):**
```bash
curl -X POST http://localhost:8080/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "amount": 100.50
  }'
```

**Withdraw (requires JWT):**
```bash
curl -X POST http://localhost:8080/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "amount": 100.50
  }'
```

**Transfer (requires JWT):**
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "amount": 50.00,
    "receiver_id": 1
  }'
```

**Get All Transactions (requires admin JWT):**
```bash
curl -X GET http://localhost:8080/transactions \
  -H "Authorization: Bearer <your-admin-jwt-token>"
```

## Development

### Running Locally (without Docker)

1. **Start PostgreSQL locally**
   ```bash
   # Make sure PostgreSQL is running on localhost:5432
   ```

2. **Update environment variables**
   ```bash
   # In your .env file, use:
   DATABASE_HOST=localhost
   DATABASE_PORT=5432
   ```

3. **Run the application**
   ```bash
   go build -o main
   ./main
   ```

## Database Schema

The application uses three main tables:

- **users**: Store user account information including roles
- **transactions**: Record all financial operations (deposits, withdrawals, transfers)

All banking operations (deposit, withdraw, transfer) are logged in the transactions table for complete audit trail.

### Database Migrations

Migrations are automatically applied when the application starts. Migration files are located in `db/migrations/` and are embedded in the binary.

To create a new migration:
```bash
migrate create -ext sql -dir db/migrations -seq your_migration_name
```

## Testing the API

You can test the API using tools like:
- **Postman**: Import the endpoint collection
- **curl**: Use the example requests above
- **HTTPie**: `http POST localhost:8080/signup username=test email=test@example.com password=test123`

## Project Structure

```
.
├── MockBankGo/
│   ├── db/
│   │   ├── migrations/      # Database migration files
│   │   └── db.go            # Database connection logic
│   ├── internal/
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── repositories/    # Data access layer
│   │   ├── services/        # Business logic layer
│   │   ├── models/          # User and transaction models
│   │   ├── apperrors/       # Custom error definitions
│   │   └── builders.go      # Dependency injection builders
│   ├── middleware/          # JWT auth & logging middleware
│   ├── auth/                # JWT generation and verification logic
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   └── main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | Database host | `localhost` |
| `DATABASE_PORT` | Database port | `5432` |
| `DATABASE_USER` | Database username | `postgres` |
| `DATABASE_PASSWORD` | Database password | - |
| `DATABASE_NAME` | Database name | - |
| `SSL_MODE` | PostgreSQL SSL mode | `disable` |
| `LISTEN_ADDR` | Server listen address | `:8080` |
| `JWT_SECRET` | JWT signing secret | - |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Port Already in Use
If you get a "port already in use" error:
- Make sure no other services are running on ports 8080 or 5433
- Or change the ports in `docker-compose.yml`

### Database Connection Issues
- Ensure your `.env` file has the correct database credentials
- Check that PostgreSQL container is healthy: `docker-compose logs db`

### Authentication Errors
- Make sure to include `Bearer ` prefix in Authorization header
- Check that your JWT token hasn't expired
- Ensure you're calling `/login` first to get a valid token

### Authorization Errors
- For admin-only endpoints, ensure your user has `admin` role
- Regular users cannot access `/transactions` endpoint
- Check user role in database: `SELECT username, role FROM users;`

### Transaction Errors  
- Verify account has sufficient balance for withdrawals/transfers
- Check that the recipient user exists for transfers
- Ensure amount values are positive numbers

### Migration Errors
- Check that your migration files are properly formatted
- Ensure migration files are in `db/migrations/` directory

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Alibi Ulanuly - alarmy126@gmail.com

Project Link: [https://github.com/Uranury/BankSystem](https://github.com/Uranury/BankSystem)