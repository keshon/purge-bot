# 🗑️ PurgeBot 🗑️

Welcome to **PurgeBot**! This bot helps manage and clean up your Discord server by automatically purging old messages from channels based on your specified duration. Whether you need to clear out outdated messages or keep your channels tidy, PurgeBot has you covered.

## 🚀 Features

- **Purge Old Messages**: Automatically delete messages older than a specified duration.
- **Stop Purging**: Easily stop the purging task for a channel.
- **List Purge Tasks**: Get a list of all active purge tasks in your guild.
- **Add or Remove Users/Roles**: Grant or revoke permission for specific users or roles to manage purge tasks.

## 🛠️ Setup

### 1. Prerequisites

- **Go**: Ensure you have Go installed on your system.
- **Discord Bot Token**: Create a bot on [Discord Developer Portal](https://discord.com/developers/applications) and get your token.
- **SQLite**: The bot uses SQLite for database storage.

### 2. Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/keshon/purge-bot.git
    cd purge-bot
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

- **Usage:** `@PurgeBot <duration>`
- **Example:** `@PurgeBot 30s` (purges messages older than 30 seconds)
- **Note:** `@PurgeBot` is just an example. The actual bot mention may vary based on the bot's name or configuration.

### Stop Purge Task

Stop the active purge task in the channel.

- **Usage:** `@PurgeBot stop`

### List Purge Tasks

Get a list of all channels with active purge tasks in the guild.

- **Usage:** `@PurgeBot list`

### Add User

Grant a user permission to manage purge tasks. You can use either username or user ID.

- **Usage:** `@PurgeBot adduser <username>` or `@PurgeBot adduserid <userID>`
- **Example:** `@PurgeBot adduser JohnDoe` or `@PurgeBot adduserid 339767128292982785`

### Remove User

Revoke a user's permission to manage purge tasks. You can use either username or user ID.

- **Usage:** `@PurgeBot removeuser <username>` or `@PurgeBot removeuserid <userID>`
- **Example:** `@PurgeBot removeuser JohnDoe` or `@PurgeBot removeuserid 339767128292982785`

### Add Role

Grant a role permission to manage purge tasks. You can use either role name or role ID.

- **Usage:** `@PurgeBot addrole <roleName>` or `@PurgeBot addroleid <roleID>`
- **Example:** `@PurgeBot addrole Admin` or `@PurgeBot addroleid 1274017921756172403`

### Remove Role

Revoke a role's permission to manage purge tasks. You can use either role name or role ID.

- **Usage:** `@PurgeBot removerole <roleName>` or `@PurgeBot removeroleid <roleID>`
- **Example:** `@PurgeBot removerole Admin` or `@PurgeBot removeroleid 1274017921756172403`

### List Permissions

Get a list of all users and roles registered to manage purge tasks, including their names.

- **Usage:** `@PurgeBot listpermissions`

### Help

Get detailed usage instructions and a list of available commands.

- **Usage:** `@PurgeBot help`

## ⚙️ Configuration

- **Purge Interval**: The interval at which the bot checks for messages to purge (default: 33 seconds).
- **Minimum Duration**: The minimum duration for purging tasks (default: 30 seconds).
- **Maximum Duration**: The maximum duration for purging tasks (default: 3333 days).

## 🗳️ Invite the Bot

To invite **PurgeBot** to your server, use the following invite link format:

`https://discord.com/oauth2/authorize?client_id=YOUR_APPLICATION_ID&scope=bot&permissions=75776`

**Required Permissions:**
- **Read Messages**
- **Send Messages**
- **Manage Messages** (for purging messages)
- **Read Message History**

Replace `YOUR_APPLICATION_ID` in the URL with your bot's actual application ID from the Discord Developer Portal.

## 📝 Example

Here's how you can use PurgeBot in your server:

1. **Start purging messages older than 1 hour:**

    ```markdown
    @PurgeBot 1h
    ```

2. **Stop purging in a channel:**

    ```markdown
    @PurgeBot stop
    ```

3. **Get a list of all purge tasks:**

    ```markdown
    @PurgeBot list
    ```

4. **Add a user to manage purge tasks:**

    ```markdown
    @PurgeBot adduser JohnDoe
    ```

5. **Remove a user from managing purge tasks:**

    ```markdown
    @PurgeBot removeuser JohnDoe
    ```

6. **Add a role to manage purge tasks:**

    ```markdown
    @PurgeBot addrole Admin
    ```

7. **Remove a role from managing purge tasks:**

    ```markdown
    @PurgeBot removerole Admin
    ```

8. **Get a list of all registered users and roles:**

    ```markdown
    @PurgeBot listpermissions
    ```

9. **Get help:**

    ```markdown
    @PurgeBot help
    ```

## 🙏 Acknowledgements

**PurgeBot** was inspired by the original [KMS Bot](https://github.com/internetisgone/kms-bot) project. The original bot, written in Python, provided the foundational concept for this Go implementation. A special thanks to the creator of that project!
