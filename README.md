# Pokémon Pick Six App

Welcome to the Pokémon Pick Six App! This application allows users to display a random Pokémon TCG card and save their favorite six cards for 24 hours. Come back tomorrow and pick a new team!

## Features

- Display a random Pokémon card.
- Save your favorite six cards.
- Cards are saved for 24 hours.

## Technologies Used

- **Frontend**: Angular
- **Backend**: Go
- **Database**: (Specify if applicable)
- **Deployment**: Fly.io

## Getting Started

### Prerequisites

- Node.js (for the frontend)
- Go (for the backend)
- Docker (optional, for containerization)

### Frontend Setup

1. Navigate to the `app` directory:
   ```bash
   cd app
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the frontend:
   ```bash
   ng serve
   ```

### Backend Setup

1. Navigate to the `server` directory:
   ```bash
   cd server
   ```

2. Build Go App:
   ```bash
   go build
   ```

3. Run the Go App:
   ```bash
   go run main.go
   ```


## Running with Docker

1. Ensure you have Docker installed on your system.
2. We will need to the Front and Backend images to run the app.
3. Navigate to the `app` directory:
   ```bash
   cd app
   ```
4. Build the Frontend Docker image:
   ```bash
   docker build -t frontend-image .
   ```

5. Navigate to the `server` directory:
   ```bash
   cd server
   ```
6. Build the Backend Docker image:
   ```bash
   docker build -t backend-image .
   ```

7. Run the Docker containers:
   ```bash
   docker run -p 3000:3000 frontend-image
   docker run -p 4200:4200 backend-image
   ```

8. Access the app at http://localhost:4200