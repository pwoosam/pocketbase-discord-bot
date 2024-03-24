# Discord Bot built using PocketBase and discordgo

## Features

- PocketBase Features
- Rock, Paper, Scissors Discord bot game

## Setup

1. Clone the repository
1. Get Discort Bot Token from [Discord Developer Portal](https://discord.com/developers/applications)
1. Invite your bot to your server using the link from the OAuth2 section of the Discord Developer Portal with the `bot` scope and the `Send Messages` permissions.
1. Create a `.env` file in the `backend/` directory of the project and add the following content:
    ```env
    DISCORD_BOT_TOKEN=<YOUR_DISCORD_BOT_TOKEN>
    ```

## How to Run

### Using Docker

1. Run the following command in the `backend/` directory of the project, run `docker-compose up`

### Using Go

1. Install `gow` by running `go install github.com/mitranim/gow@latest`
1. Run the following command in the `backend/` directory of the project, run `./start_local.sh`
