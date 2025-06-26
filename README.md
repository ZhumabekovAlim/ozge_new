# Clean Mobile App (Golang)

A clean architecture template for a mobile application backend using **Go 1.23**. This repository provides a structured approach to build scalable and maintainable backend systems.

---

## 📂 Project Structure

```
clean_mobile_app/
├── cmd/
│   └── web/
│       ├── helpers.go         # Utility functions for the web layer
│       ├── initializer.go     # Application initialization logic
│       ├── main.go            # Entry point of the application
│       ├── middleware.go      # HTTP middlewares
│       └── routes.go          # HTTP route definitions
├── config/
│   └── config.yaml            # YAML file for application configuration
├── db/
│   └── migrations/            # Database migration files
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration loading and handling
│   ├── handlers/              # HTTP request handlers
│   ├── models/                # Application data models
│   ├── repositories/          # Data access layer (DB interactions)
│   └── services/              # Business logic and services
├── go.mod                     # Go module definition
```

---

## 🚀 Features

- **Clean Architecture**: Well-structured layers for separation of concerns.
- **Go 1.23**: Latest version of Go with performance improvements.
- **Scalability**: Designed to handle complex business logic with ease.
- **Config Management**: Centralized and manageable application settings.
- **Database Migrations**: Organized schema changes under the `migrations/` directory.

---

## 📑 Contracts & Signatures API

* `POST /contracts/with-fields` – create a contract together with its additional fields. The request is a multipart form where the `fields` parameter contains a JSON array describing field names and types. The payload also accepts a `company_sign` flag (`1` for signed, `0` for unsigned).
* `GET /contracts/token/{token}/details` – retrieve a contract by token along with its signing method, `company_sign` status and list of additional fields.
* `POST /signatures` – create a signature for a contract. Along with standard fields you can pass a `field_values` JSON array to store values for the contract's additional fields in one request.

## Company API
* `POST /companies/{id}/reset-password` – set a new password for a company without providing the old one. The request body is `{ "new_password": "..." }`.


---

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ZhumabekovAlim/clean_mobile_app.git
   cd clean_mobile_app
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment**:
   - Update `config/config.yaml` with your application settings.

4. **Run the application**:
   ```bash
   go run cmd/web/main.go
   ```

---

## 🔧 Configuration

The configuration file (`config/config.yaml`) includes parameters like database credentials, server ports, etc. Modify it to suit your environment.

---

## 🗃️ Database Migrations

Use a tool like [golang-migrate](https://github.com/golang-migrate/migrate) to manage your database migrations. Place your `.sql` files in the `migrations/` folder.

To run migrations:
```bash
migrate -path migrations -database "your-database-url" up
```

---

## 🐳 Docker Deployment

1. Create a `.env` file with your database credentials (see the provided example).
2. Build and start all services with `docker-compose` which will run migrations using a dedicated container:
   ```bash
   docker-compose up --build
   ```
3. The application will be available on [http://localhost:4000](http://localhost:4000).


---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a pull request.


---

## 👨‍💻 Author

- **Alim Zhumabekov** - [ZhumabekovAlim](https://github.com/ZhumabekovAlim)

