# Gitea account registration Bot for Discord

This repository contains the source code for a Discord bot that interacts with Gitea to allow user self-registration from
within Discord. It's designed to work across multiple platforms, including Docker, Windows, Linux, and macOS.

## Description

This Discord bot is designed to react to commands on Discord and interact with the Gitea API to create and manage user
accounts. It currently supports the `!register` command to facilitate user registration through Discord.

## Environment Variables

The bot requires the following environment variables to be set:

- `DISCORD_TOKEN`: The token for your Discord bot.
- `GITEA_TOKEN`: The API token for accessing Gitea.
- `GITEA_URL`: The base URL of your Gitea instance.

### How to get the `DISCORD_TOKEN`

1. Log in to the Discord Developer Portal

2. **Create a New Application:**  
    Click on the "New Application" button. Name your application and create it.

3. **Create a Bot User:**  
    In your application's settings, navigate to the "Bot" tab.  
    Click on the "Add Bot" button and confirm the creation. This will create a new bot user associated with your application.

4. **Get the Bot Token:**  
    Under the bot settings, you will find a section labeled "TOKEN" with a link to "Copy".
    Click on "Copy" to copy the bot token. You'll use this token in your bot's configuration.

5. **Invite the Bot to Your Server:**  
    In the application settings, navigate to the "OAuth2" tab.
    Under "SCOPES", check bot.
    In the "BOT PERMISSIONS" section, select the permissions your bot requires.
    Copy the generated URL under "SCOPES" and paste it into your web browser. This will allow you to add the bot to one of your Discord servers.

### How to get the Gitea Token

1. Login to your gitea instance with a user that has access to the Administration or create a user with "Administrator" privileges.

2. Navigate to your user settings and select "Applications" in the sidebar (`https://GITEA_URL/user/settings/applications`). 
   Now you will need create an Access Token with write privileges to the `admin` scope.

3. Record the Token that is displayed when saving the token.

## Build it yourself

To build the bot from source, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://git.beisel.it/beisel.org/gitea-register-account-bot.git
   ```

2. Navigate to the cloned directory:

   ```bash
   cd gitea-register-account-bot
   ```

3. Build the binary (ensure you have Go installed):

   ```bash
   go build -o gitea-register-account-bot
   ```

### Creating a Docker Image

You may also create a Docker image using the provided Dockerfile:

1. Build the Docker image:

    ```bash
    docker build -t gitea-register-account-bot:latest .
    ```

2. Run the Docker container:

    ```bash
    docker run -d --name gitea-register-account-bot -e DISCORD_TOKEN=your_token -e GITEA_TOKEN=your_token -e GITEA_URL=your_url gitea-register-account-bot:latest
    ```

## Using the Prebuilt Images and binaries

We provide versioned as well as nightly builds through our CI/CD pipeline, where `:nightly` refers to an automatic build from
the latest commit on the `main` branch. The latest release-version can be found with the `latest` Tag.

### Using the Prebuilt Docker Image

To use the prebuilt Docker image:

1. Pull the `latest` image from the registry:

    ```bash
    docker pull git.beisel.it/beisel.org/gitea-register-account-bot:latest
    ```

2. Run the Docker container:

    ```bash
    docker run -d --name gitea-register-account-bot -e DISCORD_TOKEN=your_token -e GITEA_TOKEN=your_token -e GITEA_URL=your_url git.beisel.it/beisel.org/gitea-register-account-bot:latest
    ```

### Running Prebuilt Binaries

Prebuilt binaries for Windows, Linux, and macOS are available in the releases section of this repository.

To run the binary:

1. Download the appropriate binary for your system.
2. Make it executable (Linux/macOS):

    ```bash
    chmod +x gitea-register-account-bot-nightly-[linux-amd64,darwin-amd64]
    ```

3. Run the binary:

   ```bash
   DISCORD_TOKEN=your_discord_token GITEA_TOKEN=your_gitea_token GITEA_URL=your_gitea_url ./gitea-register-account-bot-nightly-[linux-amd64,darwin-amd64]
   ```

For Windows, simply execute the .exe file after setting your environment variables.

## Creating Gitea users from Discord

Once you have the bot up and running with the required environment variables set and invited it to your Discord server you can start creating users. The bot will listen to the `!register` command and expects the following syntax:

```bash
!register --username=[Desired Username] --email=[Your Email address]
```

The bot will open a conversation via DM with the user who issued the command and confirm whether the registration was successful with the URL of the Gitea Server, and a temporary password which has to be changed on first login.

## Contributing

Contributions to this project are welcome. Please ensure you follow the existing code style.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the repository, or hit me up on the Gitea discord @pacey.
