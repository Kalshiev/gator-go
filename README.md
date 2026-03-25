# Gator - Blog Aggregator
This is a learning project used to learn about Postgres, SQL and Go.
- This is a guided project that is part of the Backend Developer Path on boot.dev

## Requirements
1. Postgres Database
2. Go Lang

## Installation
1. Git clone this repo
```bash
$ git clone https://github.com/Kalshiev/gator-go
```
2. Install gator
```bash
$ cd gator-go
$ go install
```
3. Generate the config file in the $HOME directory
```bash
$ touch $HOME/.gatorconfig.json
$ nano $HOME/.gatorconfig.json
```
4. Config file structure
```json
{"db_url":"postgres://postgres@localhost:5432/gator?sslmode=disable", "current_user_name": ""}
```
## Usage
1. Register
```bash
$ gator register [username]
```
Creates a new user with username [username] and [password]

2. Login
```bash
$ gator login [username]
```
Logs in an existing user

3. Add Feed
```bash
$ gator addfeed [name] [url]
```
Adds a new feed with [name] and [url]

4. Aggregate
```bash
$ gator agg [timeinterval]
```
Fetches posts from feeds every [time interval] e.g. 5s, 5m, 5h

5. Browse feed
```bash
$ gator browse [limit(optional)]
```
Displays the posts from followed feeds by the current user

6. Help
```bash
$ gator help [command(optional)]
```
Displays a list for all registered commands, usage and description

## TODO
- [ ] Sorting and Filtering options to the Browse command
- [ ] Add concurrency to the agg command
- [ ] Add a search command that allows fuzzy searching of posts
- [ ] Add Bookmarking and Liking posts
- [ ] Implement a service manager that keeps agg running in the background