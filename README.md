# Medical Center Microservices

This is a microservices-based medical center application that allows administrators to manage departments, doctors to register and manage their schedules, and patients to book appointments.

## Architecture

The application is split into the following microservices:

1. **Auth Service** (Port 8081)
   - Handles user authentication and authorization
   - Manages user registration and login
   - Supports admin and doctor roles

2. **Department Service** (Port 8082)
   - Manages medical departments
   - Allows administrators to create and update departments

3. **Doctor Service** (Port 8083)
   - Manages doctor profiles
   - Handles doctor registration
   - Manages doctor availability and schedules

4. **Appointment Service** (Port 8084)
   - Handles appointment booking
   - Manages appointment status
   - Allows patients to book appointments without registration

5. **API Gateway** (Port 8080)
   - Routes requests to appropriate microservices
   - Handles request/response transformation
   - Provides a unified API interface

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd medical-center
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. The services will be available at:
   - API Gateway: http://localhost:8080
   - Auth Service: http://localhost:8081
   - Department Service: http://localhost:8082
   - Doctor Service: http://localhost:8083
   - Appointment Service: http://localhost:8084

## API Endpoints

### Authentication
- POST /api/v1/auth/register - Register a new user (admin/doctor)
- POST /api/v1/auth/login - Login user

### Departments
- POST /api/v1/departments - Create a new department (admin only)
- GET /api/v1/departments - List all departments
- GET /api/v1/departments/:id - Get department details
- PUT /api/v1/departments/:id - Update department (admin only)

### Doctors
- POST /api/v1/doctors - Register a new doctor (admin only)
- GET /api/v1/doctors - List all doctors
- GET /api/v1/doctors/:id - Get doctor details
- PUT /api/v1/doctors/:id - Update doctor (admin only)
- POST /api/v1/doctors/:id/availability - Set doctor availability
- GET /api/v1/doctors/:id/availability - Get doctor availability

### Appointments
- POST /api/v1/appointments - Book a new appointment
- GET /api/v1/appointments - List all appointments
- GET /api/v1/appointments/:id - Get appointment details
- GET /api/v1/appointments/doctor/:doctor_id - Get doctor's appointments
- GET /api/v1/appointments/department/:department_id - Get department's appointments
- PUT /api/v1/appointments/:id - Update appointment status

## Default Admin Account

A default admin account is created when the system starts:
- Email: admin@example.com
- Password: admin123

## Development

To run the services locally for development:

1. Install dependencies for each service:
   ```bash
   cd services/<service-name>
   go mod download
   ```

2. Run each service:
   ```bash
   go run main.go
   ```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 