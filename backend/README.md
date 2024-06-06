# Coding Blog

A place for me to document things I have learned.

## Development
### Run locally
```bash
./scripts/run.sh
```

## Code Architecture.
### Mostly follows **clean code architecture**
- Layers
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
    - **Handlers/Usecase**
        - Core app logics, uses repository layer for CRUD operations
- All layers communicate through interfaces
- Independent of database
- There is NO Controller/Delivery layer to abstract different input methods (ex: grpc)
    - `handlers` depands on `net/http`

### Architecture diagram
TODO

### Class diagram
TODO


## TODO
- Remove unassasary interfaces.... 
    - Preferably, only use interface for connecting handlers and repositories (will change DB, or use ORM)
- Better naming. EX: Models -> Tables
