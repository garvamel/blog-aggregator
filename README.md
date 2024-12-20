# Aggre-Gator

## Description

Gator is a small command-line RSS feed aggregator that allows different users to follow feeds.

## Setup

To run gator you need to set up the _go toolchain_ a _postgresql_ database and the _goose_ tool.
You can check how to install _go_ and _postresql_ [here](https://webinstall.dev/golang/) and [here](https://www.boot.dev/lessons/74bea1f2-19cd-4ea9-966e-e2ca9dd1dfa9), respectively

## Config

Gator needs config file in your user folder with the following (prettyfied) format:

```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

When filled with correct information, it should look something like this:

```json
{"db_url":"postgres://giapoldo:@localhost:5432/gator?sslmode=disable","current_user_name":"Giapoldo"}
```

Save it in your home folder as `.gatorconfig.json` (the initial dot is necessary)

## Install

1. Install _goose_ by running `go install github.com/pressly/goose/v3/cmd/goose@latest

Install Gator by running `go install github.com/giapoldo/blog-aggregator``

## Usage

you can run the following commands:

### To manage users

1. `gator register <username>`
2. `gator login <username>`

### To manage RSS feeds

1. `gator addFeed "<feed name>" "<feed url>"`
2. `gator follow "<feed url>`
3. `gator following`

## To browse recent posts

1. `gator browse [optional limit]`

## To reset the RSS feed database

1. `gator reset`