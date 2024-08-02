# 🗑️ PurgeBot 🗑️

Welcome to **PurgeBot**! This bot helps manage and clean up your Discord server by automatically purging old messages from channels based on your specified duration. Whether you need to clear out outdated messages or keep your channels tidy, PurgeBot has you covered.

## 🚀 Features

- **Purge Old Messages**: Automatically delete messages older than a specified duration.
- **Stop Purging**: Easily stop the purging task for a channel.
- **List Purge Tasks**: Get a list of all active purge tasks in your guild.
- **Help Command**: Get a comprehensive list of commands and usage instructions.

## 🛠️ Setup

### 1. Prerequisites

- **Go**: Ensure you have Go installed on your system.
- **Discord Bot Token**: Create a bot on [Discord Developer Portal](https://discord.com/developers/applications) and get your token.
- **SQLite**: The bot uses SQLite for database storage.

### 2. Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/keshon/purge-bot.git
    cd your-repo
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up your environment:**

    Create a `.env` file in the root directory with the following content:

    ```env
    DISCORD_KEY=your-discord-bot-token
    ```

4. **Run the bot:**

    ```bash
    go run main.go
    ```

## 📜 Commands

### Purge Old Messages

Automatically purge old messages in the channel.

- **Usage:** `@bot <duration>`
- **Example:** `@bot 30s` (purges messages older than 30 seconds)

### Stop Purge Task

Stop the active purge task in the channel.

- **Usage:** `@bot stop`

### List Purge Tasks

Get a list of all channels with active purge tasks in the guild.

- **Usage:** `@bot list`

### Help

Get detailed usage instructions and a list of available commands.

- **Usage:** `@bot help`

## ⚙️ Configuration

- **Purge Interval**: The interval at which the bot checks for messages to purge (default: 33 seconds).
- **Minimum Duration**: The minimum duration for purging tasks (default: 30 seconds).
- **Maximum Duration**: The maximum duration for purging tasks (default: 3333 days).

## 🗳️ Invite the Bot

To invite **PurgeBot** to your server, use the following invite link format:

`https://discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID&scope=bot&permissions=75776`

**Required Permissions:**
- **Read Messages**
- **Send Messages**
- **Manage Messages** (for purging messages)
- **Read Message History**

Replace `YOUR_CLIENT_ID` in the URL with your bot's actual client ID from the Discord Developer Portal.

## 📝 Example

Here's how you can use PurgeBot in your server:

1. **Start purging messages older than 1 hour:**

    ```markdown
    @bot 1h
    ```

2. **Stop purging in a channel:**

    ```markdown
    @bot stop
    ```

3. **Get a list of all purge tasks:**

    ```markdown
    @bot list
    ```

4. **Get help:**

    ```markdown
    @bot help
    ```

## 🙏 Acknowledgements

**PurgeBot** was inspired by the original [KMS Bot](https://github.com/internetisgone/kms-bot) project. The original bot, written in Python, provided the foundational concept for this Go implementation. A special thanks to the creator of that project!