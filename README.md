<p align="center">
<a href="" target="">
<img src="https://raw.githubusercontent.com/jyotirmoydotdev/Openfy/main/src/OpenfyLogo.svg" alt="Openfy logo">
</a>
</p>

> ⚠️ Note: This project is currently under development and may not be fully functional. Feel free to explore the code, but be aware that some features may not work as intended.

# Openfy

Openfy is a self-hosted e-commerce platform built with Go and Gin, providing a flexible and customizable solution for managing products and handling customer authentication.

## Features

- **Product Management**: Easily create, update, and delete products.
- **Customer Authentication**: Secure customer authentication with JWT tokens.
- **StaffMember Panel**: StaffMember-specific routes for managing products.
- **Flexible and Extendable**: Built with Go and Gin for flexibility and extensibility.

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/jyotirmoydotdev/openfy.git
   cd openfy
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up your configuration:

   Copy the `.env.example` file to `.env` and configure your environment variables, such as database connection details and secret keys.

4. Run the application:

   ```bash
   go run main.go
   ```

   The application will be accessible at `http://localhost:8080`.

## Routes

- **Customer Routes**:
  - `POST /signup`: Register a new customer.
  - `POST /login`: Customer login.
  - `GET /products`: Retrieve details of all active products.

- **StaffMember Routes**:
  - `POST /staffMember/signup`: Register a new staffMember customer.
  - `POST /staffMember/login`: StaffMember login.


- **Authenticated StaffMember Routes** (Protected by JWT):
  - `GET /staffMember/product`: Retrieve details of a specific product.
  - `GET /staffMember/products`: Retrieve details of all products.
  - `POST /staffMember/product/new`: Create a new product.
  - `PUT /staffMember/product`: Update an existing product.
  - `DELETE /staffMember/product`: Delete a product.
  - `DELETE /staffMember/variant`: Delete a product variant.

Explore the features, test functionalities, and provide us with valuable feedback. Your insights are crucial in shaping the final release, and we appreciate your participation in this preview phase.
