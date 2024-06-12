# Coding Blog

A place for me to document things I have learned.

## Development
### Run locally
```bash
./scripts/run.sh
```

## Tech stack
### Backend
Mostily uses golang's builtin librariy
- **net/http** for server
- **http.ServeMux** for routing
- **database/sql** for querying databases
- **goose** for database migrations
- Supported Databases: 
    - SQLite
    - PostgreSQL [TODO]

### Frontend [TODO]
- htmx

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


## TODO
### High
- tests:
    - repository intergration tests [doing]
    - auth related unit tests 
    - handler unit tests

### Mid
- Custom error returns
    - checks input and group them in one response

### Low
- Log request after response (with status code)
