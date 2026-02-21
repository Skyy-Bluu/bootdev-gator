# Gator- An RSS feed aggregator

## Description

A CLI tool that allows the user to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Dependencies

* ### Golang 1.25+
  
  Install on MacOS/Linux:
  ```
  curl -sS https://webi.sh/golang | sh; \
  source ~/.config/envman/PATH.env
  ```
  Install on Windows:
  ```
  curl.exe https://webi.ms/golang | powershell
  ```

  Ensure you installation worked with
  ```
  go version
  ```
  
* ### Postgresql 15
  
  Install on MacOS using [brew](https://brew.sh/):
  ```
  brew install postgresql@15
  ```

  Linux / WSL (Debian). Here are the [docs from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql), but simply:
  ```
  sudo apt update
  sudo apt install postgresql postgresql-contrib
  ```
  
  Ensure you installation worked with
  ```
  psql --version
  ```

  (Linux only) Update postgres password:
  ```
  sudo passwd <PASSWORD>
  ```
  Replace placeholder text with a password you'll remember

  Enter the psql shell:

  Mac:
  ```
  psql postgres
  ```
  Linux:
  ```
  sudo -u postgres psql
  ```
  You should see
  ```
  postgres=#
  ```
  Create a new database:
  ```
  CREATE DATABASE <DB NAME>;
  ```
  For Linux only, connect to the db you created:
  ```
  \c <DB NAME>
  ```
  Set the user password
  ```
  ALTER USER postgres PASSWORD '<PASSWORD>';
  ```
  Run the migration scripts in order
  ```
  psql -d [DB NAME HERE] -f sql/schema/001_users.sql
  psql -d [DB NAME HERE] -f sql/schema/002_feeds.sql
  psql -d [DB NAME HERE] -f sql/schema/003_feed_follows.sql
  psql -d [DB NAME HERE] -f sql/schema/004_feeds_add_last_fetched_at.sql
  psql -d [DB NAME HERE] -f sql/schema/005_posts.sql
  ```

  - ### Create congiuration file
    Create ```.gatorconfig.json``` in your home directory adding the you database URL replacing in your username, password and database name. For MacOS, leave the password empty
    ```
    {
    "db_url": "postgres://<username>:<password>@localhost:5432/<db name>"
    }
    ```
### Installing

```
go install bootdev-gator
```

### Executing program
Here is a list of the commands available and description


| Command | Description |
| ------------- | ------------- |
| register \<username\> | Create a new user |
| login \<username\> | Login as user |
| users | List all users |
| reset | Delete all records in all tables |
| addfeed \<name\> \<url\> | Add new RSS feed  |
| feeds | List all feed |
| follow \<url\> | Follow feed |
| unfollow \<url\> | Unfollow feed |
| following | List followed feeds for logged in user |
| agg \<interval\> | Continuously fetch posts from the  added RSS feeds |
| browse \<limit\> | View posts from the logged in user's followed RSS feeds. Default limit is 2 |


## Acknowledgments

* Lane Wagner - [boot.dev](boot.dev)
