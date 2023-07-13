# API Routes Documentation

This application provides a RESTful API for managing users, wallets, subscriptions, phonebook contacts, and SMS services. The following API endpoints are available:

### User Management

- `POST /signup`: register a new user.
- `POST /login`: log in a user and return a JWT token.
- `GET /logout`: log out a user and invalidate the JWT token.

### Subscription Management

- `PUT /numbers`: create and subscribe to a new number.
- `POST /numbers/buy-rent`: buy or rent a new number and subscribe.

### Wallet Management

- `POST /wallets/charge-request`: request to charge the wallet balance.
- `POST /wallets/finalize-charge`: finalize the wallet charge request.

### Phonebook and Contact Management

- `GET /phonebook`: get all phonebook contacts for the authenticated user.
- `GET /phonebook/username`: get a single phonebook contact by the username of the authenticated user.
- `DELETE /phonebook`: delete a phonebook contact.
- `POST /phonebook`: create a new phonebook contact.
- `PUT /phonebook`: edit an existing phonebook contact.
- `POST /contact`: create a new SMS contact.
- `PUT /contact`: edit an existing SMS contact.
- `GET /contact`: get all SMS contacts for the authenticated user.
- `GET /contact/phonebook`: get all SMS contacts for a phonebook.

### SMS Management

- `POST /sms`: send a new SMS message to a single recipient.
- `POST /sms/periodic`: send a new periodic SMS message to a single recipient.
- `POST /sms/username`: send a new SMS message to all contacts for a specific username.
- `POST /sms/username/periodic`: send a new periodic SMS message to all contacts for a specific username.

### SMS Template Management

- `POST /templates/new`: create a new SMS template.
- `GET /templates`: get all SMS templates for the authenticated user.
- `POST /templates/sms`: send a new SMS message using an existing template.
- `POST /templates/sms/periodic`: send a new periodic SMS message using an existing template.
- `POST /templates/sms/username`: send a new SMS message to all contacts for a specific username using an existing template.
- `POST /templates/sms/username/periodic`: send a new periodic SMS message to all contacts for a specific username using an existing template.

### Admin Management

- `GET /admin/disable-user/:userId`: disable a user.
- `GET /admin/change-pricing`: change pricing.
- `GET /admin/sms-report/:userId`: get SMS history for a specific user.

All endpoints require authentication, except for `POST /signup` and `POST /login`.

## Authentication

Authentication is handled using a JSON Web Token (JWT). The client must include a JWT token in the `Authorization` header of each API request, except for `POST /signup` and `POST /login`.

The JWT token can be obtained by logging in with a valid email and password. The token is valid for a limited time and must be renewed periodically.

## Authorization

Authorization is handled using custom middleware that verifies the user's role and permissions for each API request. Some endpoints are only accessible to admin users.

## Error Handling

The API returns appropriate HTTP status codes for each request. In case of an error, the response body will contain an error message in JSON format.