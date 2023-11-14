# Openfy 

Openfy is a self-hosted e-commerce platform built with Go and Gin, providing a flexible and customizable solution for managing products and handling user authentication.

## Features

- **Product Management**: Easily create, update, and delete products.
- **User Authentication**: Secure user authentication with JWT tokens.
- **Admin Panel**: Admin-specific routes for managing products.
- **Flexible and Extendable**: Built with Go and Gin for flexibility and extensibility.

## Prerequisites

Before running Openfy, ensure you have the following installed:

- Go (1.15 or higher)
- Gin (Golang web framework)

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

- **Public Routes**:
  - `GET /products`: Retrieve a list of products.
  - `GET /products/:id`: Retrieve details of a specific product.

- **Admin Routes** (Protected by JWT):
  - `POST /product/new`: Create a new product.
  - `PUT /product/:id`: Update an existing product.
  - `DELETE /product/:id`: Delete a product.
  - `GET /admin`: Admin-specific endpoint.

## Authentication

- User authentication is implemented using JWT tokens.
- To authenticate, include the JWT token in the "Authorization" header using the "Bearer" scheme.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Make sure to replace placeholder details like `yourusername` and update the project-specific sections based on your actual project structure, features, and configurations. Include any additional information that might be relevant to users and contributors.