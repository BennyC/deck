# Deck API Service

## Requirements

* Go 1.15.6


## Description

A REST API to simulate a deck cards.

This has been created with an in-memory storage provider and interfaced services
to perform the operations required for this API to function.

As the storage provided is interfaced, this could be quickly swapped out for a true database.

Integration tests have been created to run against the API endpoints. Basic unit tests have been included for both services.

Simple middlewares have also been included for logging and automatically setting the content type header.

A Postman collection has also been included to test endpoints with

## Getting started

Run integration tests:

```make integration```

Run unit tests:

```make unit```

Run __all__ tests:

```make test```

Run service without building:

```make run```

Build service:

```make build```

## Extra Information

A `PORT` environment variable can used to run the service on a different port, will default to `8080`.

Example: `PORT=3001 make run`