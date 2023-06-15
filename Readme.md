# Blog API

This is a RESTful service built using Golang that provides CRUD (Create, Read, Update, Delete) operations for a blogging application. The API includes JWT (JSON Web Token) authentication for secure access and utilizes MongoDB as the underlying database.

## Motivation

The primary motivation behind this project is to enhance proficiency in Golang programming language and gain hands-on experience with integrating MongoDB and JWT authentication into an application. By building this Blog API, the aim is to improve knowledge and skills in these areas and provide a solid foundation for future projects.

## Features

- Create, read, update, and delete blog posts
- User registration and authentication using JWT tokens
- Secure endpoints requiring valid authentication
- Persistent storage of blog posts in MongoDB
- RESTful API design principles followed

## Deployment

The project is currently hosted on [https://blog-api-d009.onrender.com](https://blog-api-d009.onrender.com) and is accessible for testing and integration purposes. Please refer to the API documentation below for detailed information on the available endpoints and their functionalities.

## API Documentation

The API documentation is available on Postman. You can find the collection and detailed descriptions of each endpoint [here](https://galactic-equinox-827112.postman.co/workspace/My-Workspace~29d281b0-81dc-440e-8d9c-0246e619554f/api/86ead33d-d7b4-4a1e-99b0-c776756ea2a9).

## Technologies Used

- Golang
- MongoDB
- JSON Web Tokens (JWT)

## Installation and Setup

To set up the project locally, follow these steps:

1. Clone the repository: `git clone https://github.com/A-Victory/Blog-API.git`
2. Navigate to the project directory: `cd Blog-API`
3. Install the required dependencies: `go mod tidy`
4. Build and run the application: `go run main.go`

## Database

The API uses MongoDB as the database for storing blog posts and user information. Make sure you have MongoDB installed and running locally or provide the connection details in the `.env` file.

Make sure you have Golang and MongoDB installed on your system before proceeding with the installation.

## Contributions

Contributions are welcome! If you'd like to contribute to the Blog-API project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make the necessary changes and commit your code.
4. Push your changes to your fork.
5. Submit a pull request to the main repository.

