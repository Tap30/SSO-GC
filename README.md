# SSO-GC API Documentation

## Overview

The SSO-GC (Single Sign-On for Go Applications) project provides a robust implementation of OAuth 2.0 and OpenID Connect protocols, facilitating secure authentication and authorization functionalities across distributed systems. This system manages user sessions and tokens, enabling applications to authenticate users and access their profile information securely.

## Components

### `config`
Handles application configuration loaded from YAML files using Viper. It manages configurations such as SSO issuer URL, server port, client ID, and client secret.

### `handler`
Manages HTTP request handling, interfacing with the authentication logic to serve OpenID Connect configuration, tokens, user information, and logout functionalities.

### `auth`
Contains core authentication logic, including the retrieval of OpenID configurations, token generation, user information fetching, and custom claim handling.

### `server`
Sets up the HTTP server using Echo, defining routes and middleware configurations necessary for handling OAuth and OpenID Connect requests.

## APIs

### OpenID Configuration

**Endpoint:** `GET /.well-known/openid-configuration`

Retrieves the OpenID Connect configuration which includes endpoints and capabilities of the OpenID provider.

### Token Generation

**Endpoint:** `POST /token`

Generates tokens based on request parameters. This can include:

- **Access Token:** Used to access protected resources.
- **ID Token:** Contains user profile information in a JWT format.
- **Refresh Token:** Used to renew access tokens without user interaction.

**Parameters:**

- `code`: Authorization code received during user authentication.
- `redirect_uri`: URI to redirect users after authentication.
- `grant_type`: Specifies the type of token request (e.g., authorization_code, refresh_token).
- `client_id`: Registered client identifier.
- `client_secret`: Secret used to authenticate the client to the token endpoint.

### User Information

**Endpoint:** `POST /userinfo`

Retrieves user information using the access token provided in the Authorization header. This endpoint decodes the access token to fetch user attributes.

### Logout

**Endpoint:** `GET /logout`

Terminates the user's session and clears relevant cookies. Optionally redirects the user to a specified URI after logout.

## Security and Compliance

This project implements standard security protocols and complies with OAuth 2.0 and OpenID Connect specifications to ensure secure transmission of information. CORS is configured for cross-origin resource sharing, allowing the server to interact securely with resources from different domains.

## Configuration

The application's configuration is managed through a `config.yaml` file, which must specify keys such as `server_port`, `sso_issuer`, `client_id`, and `client_secret`. Errors in configuration loading or parsing are handled with immediate log output and system halt to prevent startup with incorrect settings.
