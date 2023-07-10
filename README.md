# Go SMS

This Go application is a backend implementation for a phonebook and SMS service application. It provides a range of features such as user authentication, number subscription, SMS service, SMS templating, Periodic SMS send and more.

## Getting Started

Clone the repository:
`git clone https://github.com/the-go-dragons/final-project2.git`

To run this application you have two aptions:

1. First is to run it using docker-compose for production :

    ```shell
    docker-compose -f docker-compose-prod.yml up -d
    ```
2. Second is to run just the dependecy in docker-compose for development :

	```shell
    docker-compose up -d
    ```
    And the run the go:

    ```shell
    go run ./cmd
    ```

## Features

### This application provides a range of features including:

- #### User authentication: users can sign up, log in and log out.
- #### Subscription management: users can subscribe to numbers.
- #### Wallet management: users can charge their wallet and manage wallet balance.
- #### SMS service: users can send SMS messages to phonebook contacts or just a number.
- #### SMS Templating: users can create a sms template and send sms by using that template.
- #### Periodic SMS: users can create a periodic sms send by setting a repeatation count and period time.
- #### Combination all features: users also can combine send sms to phone books or to a number with templating and period send.
- #### Admin system: admins can disable user, get a sms report, set sms filter by word or word regex, set pricing for sms and more.

## Code Structure

- `app`: contains the main application logic for the backend.
- `middleware`: contains custom middleware used in the application.
- `interfaces/http`: contains http handlers for each of the supported endpoints.
- `interfaces/persistence`: contains implementations of the persistence interface to interact with the database.
- `usecase`: contains the business logic and use case implementation of the application.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

<p style="text-align: center; width: 100%; ">Copyright&copy; 2023 <a href="https://github.com/the-go-dragons">The Go Dragons Team</a>, Licensed under MIT</p>