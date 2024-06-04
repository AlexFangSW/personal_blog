# Coding Blog

A place for me to document things I have learned.

## Development
### Run locally
```bash
./scripts/run.sh
```

## Code Architecture.
### Follows **Clean code architecture**
- Layers
    - entities
    - models
    - repository
    - handlers
- All layers communicate through interfaces
- Independent of database
- Independent of server framework (TODO)
    - need to add another layer between `handlers` and `repository` (`usecase` layer)
    - `handlers` depands on `net/http` related frameworks

### Architecture diagram
TODO

### Class diagram
TODO

