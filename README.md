# User Management System Blueprint Server

![Blueprint](/static/logos/logo_banner.png)

## Description

This project provides a user management system for the Blueprint Server. Stevens Blueprint uses Authelia as
SSO application. Authelia provides a YAML file called users.yaml where all the authorized users are stored.
This API intends to provide a way to manage the users.yaml. The API should provide endpoints to create, add,
delete, update and disable users. The service will follow a RESTful API architecture. 

## Installation

1. Clone the repository:
```
git clone https://github.com/your-username/user-management-system.git
```

2. Build the project:

```
go build
```

3. To run the server. Use the -dev flag to run the service in dev mode.
```
go run main.go -dev
```
Note: If the server is running in dev mode, then the yaml file ```users.yaml``` in the root directory will be used.
If the server is running in prod mode, then the path to the yaml file has to be specified in a ```.env``` file

## Running the service
Run
```
go run main.go -dev
```

## Docs
Available Endpoints
``
GET /api/v1/all
``

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these guidelines:


1. Create a new branch for your feature or bug fix.

2. Make your changes and commit them with descriptive commit messages.

3. Push your changes to your branch.

4. Submit a pull request, explaining the changes you have made.

## Testing

To run the tests for this project, use the following command:

```
cd tests
go test
```

## License

This project is licensed under the [MIT License](LICENSE).
