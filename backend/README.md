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
Will need JWT token to use some of the APIs, such as create, update, delete and listing unvisible resources.

- Swagger Doc: [swagger.json](./docs/swagger.json)
- Blogs
    - [x] Basic CRUD operations
    - List filters
        - [x] By topic ids
        - [x] By topic and tag ids
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

## TODO
- Rate Limit (login)
