# Coding Blog Backend

A place for me to document things I have learned.

And as a opportunity to test out go 1.22's new "net/http" update.

Which adds support of limiting request method in ServieMux routing pattern.

## Development
### Run locally
```bash
./scripts/run.sh
```
### Docker
```bash
./scripts/docker_build.sh
```
```bash
./scripts/docker_run.sh
```

## Tech stack
Mostily uses golang's builtin librariy
- **net/http** for server
- **http.ServeMux** for routing
- **database/sql** for querying databases
- **goose** for database migrations
- **SQLite** as the database

## Code Architecture.
- **Entities**
    - Basic structs
- **Models**
    - Focuses on making sql queries on specific tables 
    - Including:
        - blogs
        - tags
        - topics
        - blog_tags (many to many)
        - blog_topics (many to many)
- **Repository**
    - A interface for CRUD operations on base tables such as: blogs, tags, topics
    - Automaticly maintains many-to-many tables: blog_tags, blog_topics
- **Handlers**
    - Core app logics, uses repository layer for CRUD operations

### Architecture diagram
TODO

### Database relations
TODO

### API Documentation
Will need JWT token to use some of the APIs, such as create, update, delete and listing invisible resources.

> [Swaggo](https://github.com/swaggo/swag) (auto genrate swagger.json) dosn't support JWT auth yet. 
> So the auth part is missing in the swagger doc. 😔 

- Swagger Doc: [swagger.json](./docs/swagger.json)

#### Progress
- Blogs
    - [x] Basic CRUD operations
        - List operations will return empty 'content' field to reduse size.
    - List filters
        - [x] By topic ids
        - [x] By topic and tag ids
        - [x] Option to return simple output with tags and topics as slugs (orignaly returns full struct of tags and topics)
            - This reduces the size from 1M to about 310K on 1000 blogs with 2 to 3 tags and topics
    - [x] md5 to check if content is the same.
- Tags
    - [x] Basic CRUD operations
    - List filters
        - [x] By topic id ( in relation to blogs under a specific topic )
- Topics
    - [x] Basic CRUD operations

### Tests
- repository integration test
    - blogs
        - [x] Basic CRUD
        - List filters
            - [x] By topic ids
            - [x] By topic and tag ids
    - tags
        - [x] Basic CRUD
        - List filters
            - [ ] list tags by topic id
    - topics
        - [x] Basic CRUD
- Auth util unit test
    - [x] jwt helper
    - [x] auth helper
- handler unit test
    - [ ] blogs
    - [x] tags
    - [ ] topics

## Tools
### [Sync tool](./cmd/sync-tool/main.go) [TODO]
I want to use my own editor to write notes.

This is a tool that can sync my notes to the server.

The notes should be organized like `dummyData` folder
- A **meta.yaml** containing tags and topics
- **blogs** folder containing blogs with frontmatter

### Referential integrity
If some blog's frontmetter has a tag or topic that doesn't exist,
it will be recorded, logged out and written to a file, after which the sync process will be terminated.
Only after passing the validation process will the data be synced to the server.

### [User register](./cmd/register/main.go)

This project is only used by one person, with no intention of saving other user's stuff.

And because of this, there is no api for registoring a new user.

Only someone with direct access to the database can register.

#### Functions
- CRUD for user table, directly operates on the database.

## TODO
- Rate Limit (login)
- Remove unecessary pointer return
- SQLite
    - activate WAL
    - command to manually vacuum "db" and "wal"
- Refector server run() (model handler buildup)
- Frontend... nice links (/<id>/<slug>)
