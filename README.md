# Go Image Service

This is a Go project that consists of two services. The first service is responsible for user authentication using JWT tokens. The second service is responsible for uploading images to a MongoDB database and downloading segmented images.


## Configuration

Before running the project, you will need to set up the config file `image_service_config.json` and set up environment variables, have a look at file `.env.template`

## Architecture
<img width="600" alt="Screenshot 2023-04-23 at 01 41 36" src="https://user-images.githubusercontent.com/34031791/235322269-77725c25-6e2d-4077-8a25-e5f0a53ecdf5.png">


## Installation

This project can be run using Docker Compose.

1. Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) on your system.
2. Clone this repository to your local machine.

## Usage

1. Navigate to the project root directory.
2. Run `docker-compose up` to start the services.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
