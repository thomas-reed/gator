# gator
Boot.dev project - A multi-user CLI-based RSS feed aggregator 

## Installation
Gator uses PostgreSQL DB (v18.1+) and is written in Go (v1.25.5+), so ensure both are installed.

1. Clone this repo locally
2. Create a database called `gator`
3. Save the connection string to the database in a file in your home folder called `.gatorconfig.json`:
```
{
  "db_url": "postgres://<db_user>:<db_password>@localhost:5432/gator?sslmode=disable",
}
```
4. Install the tool: `go install .`

## Usage 
`gator <command> [<args..>]`

Available commands:
* `register <username>`            - adds a user to the db and sets user as the current user in the config file
* `login <username>`               - Sets the given user as the current user in the config file
* `users`                          - Lists registered users
* `addfeed <feed_name> <feed_url>` - Adds a feed to the database, and follows the feed for the current user
* `feeds`                          - Lists all feeds that have been added to the database
* `follow <feed_url>`              - Follows the given feed URL for the current user, provided it has already been added to the database
* `following`                      - Lists all the feeds the current user is following by name
* `unfollow <feed_url>`            - Unfollows given feed for the current user
* `browse <post_limit[optional]>`  - Displays the most recent posts from your feeds.  Number of posts is set by `<post_limit>` (default 2)
* `agg <poll_interval>`            - Aggregates posts from all feeds and stores them in the DB.  `<poll_interval>` e.g. 10s, 5m, 1h, etc. how often to check feeds for new posts.
* `reset`                          - (DESTRUCTIVE) If you want to reset your database, here you go. You've been warned :)

