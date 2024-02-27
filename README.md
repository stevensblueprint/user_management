# User Management System Blueprint Server

![Blueprint](/static/logos/logo_banner.png)

## Description

Go-based User Management service is designed to interface with a users.yaml file provided by Authelia for Single Sign-On (SSO) configurations. It offers a comprehensive suite of functionalities tailored for managing user access and credentials within an SSO environment.

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
### User Schema
```
Displayname string
Email string
Password string
Disabled bool
Groups []string
```

### Available Endpoints
```
GET /api/v1/users/user?username={username}
```
Returns a user with provided username in url params

```
POST /api/v1/users/user
```
Adds user to yaml file

```
PUT /api/v1/users/user?username={username}
```
Update a user and writes updated user to yaml file

```
DELETE /api/v1/users/user?username={username}
```
Deletes a user from yaml file

```
GET /api/v1/users/all
```
Returns all users in yaml file

```
POST /api/v1/users/user/enable?username={username}
```
Sets enabled field of user given in URL param to true

```
POST /api/v1/users/user/disable?username={username}
```
Sets enabled filed of user given in URL param to false

```
POST /api/v1/users/register
```
Creates a token to register a user and adds it to pool of valid tokens. Sends invitation to the new user
to finish creating the account. 

```
PUT /api/v1/users/reset_password
```
Updates password for username found in the body of the request

```
GET /api/v1/users/health
```
Health check

## Workflow Overview
To maintain the API secure and to block unwanted users from creating an account, an admin user will have to register the user using the ```POST /api/v1/users/register``` endpoint. This endpoint will create a valid token that will be sent to the new user via email. The content of the email will be a welcome message, along with a URL to complete the registration using the ```POST /api/v1/users/user```. Users will only by able to complete the registration once, because the token sent in the welcome email will expire once the user sends the form.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these guidelines:


1. Create a new branch for your feature or bug fix.
```
git checkout -b {feature/fix}/name-of-branch
```

2. Make your changes and commit them with descriptive commit messages.
```
git add {modified files}
git commit -m "commit message"
```

3. Push your changes to your branch.
```
git push -u origin {branch name}
```

4. Submit a pull request, explaining the changes you have made. If you open the code on GitHub, you'll see a green button that says "Create Pull Request"

## Testing

To run the tests for this project, use the following command:

```
cd tests
go test
```

## License

This project is licensed under the [MIT License](LICENSE).
