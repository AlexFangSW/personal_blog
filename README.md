# Coding Notes
Things that I find useful.

[Build status XXX] 

Link to the website: [PLACE THE LINK HERE]

## Motive
A opportunity to practice **GO** while creating something I could use.

This project started around when **GO 1.22**  was released.
At that time I was chosing which framework or libraries I sould use. 
In the release notes, one thing cought my eye: 
**[Enhanced routing patterns](https://tip.golang.org/doc/go1.22#enhanced_routing_patterns)**.
This improved `net/http.ServeMux` quite a lot, adding features such as **path params**, **restrict methods to routes**...etc.

After seeing this plus a video on youtube by **Dreams of Code**: [The standard library now has all you need for advanced routing in Go](https://www.youtube.com/watch?v=H7tbjKFSg58),
I thought, 'why not give the standard library a try ?' After than, this project was created.

## How it is used
I enjoy using my own editor, and I also want to seperate my content from the source code.

- Deploy the server, and create a user with **[UserRegister](#user-register)** CLI tool.
- Write blogs and organize it like [dummyData](./backend/dummyData/) directory.
- Use **[SyncTool](#synctool)** to sync data to the server.

> Tools are mentioned in [CLI tools](#cli-tools) section.

## Project Stucture
This project is seperated into **frontend**, **backend** and **CLI tools** for me to sync my data to the server.

> I had thought of using **[Htmx](https://htmx.org/)**, **[Templ](https://templ.guide/)** and **GO**,
but I was already half done with the backend, so that will be for another time.

## Backend
This part is mostly made with the **GO** standard library,
not just routing, but also database related stuff.

> A one point I read about [Clean code architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
and refactored the backend to strictly follow **clean code**. After which I **refactored part of it** again,
finding that by strictly following **clean code** caused the project to feel overengineered and filled with unnecessary interfaces.

### Tech stack
- Language: **GO**
- Server: **net/http** ( Routing is done with `net/http.ServeMux`)
- ORM: **None**, this project uses `database/sql` with **raw sql queries**.
- Database: **SQLite**
- Database Migration: **[goose](https://github.com/pressly/goose)**
- API documentation: generated with **[swaggo](https://github.com/swaggo/swag)**

> SQLite [WAL](https://www.sqlite.org/wal.html) mode enables **none blocking**
read and writes.

### API Documentation
#### Overview
APIs are seperated into **PUBLIC** and **PRIVATE**, 
**PUBLIC** APIs can be access by anyone, while **PRIVATE** APIs 
needs **JWT** token.

**JWT** is also stored in the database, in order to use **PRIVATE** APIs,
one must have a valid **JWT** token, while also matching the entry in the database.

-   <details>
    <summary>Blogs API</summary>

    - **Public API** ( Access blogs that are visible or not soft deleted):
        - List
            - all
            - filter by topic id (allow multiple ids)
            - filter by topic and tag ids (allow multiple ids) 
        - Get by id
    - **Private API** ( Needs JWT token, have access to all blogs regarding visibility or soft delete status )
        - Create
            - auto generate id
            - with specified id
        - List
            - all
            - filter by topic id (allow multiple ids)
            - filter by topic and tag ids (allow multiple ids) 
            - simplified
                - only includes necessary fields to verify change, such as: **content_md5**, **tag.slugs**, **topic.slugs**...etc.
                  (used by **SyncTool**)
        - Update
        - Delete
            - soft delete
            - restore soft deleted blog
            - delete

    </details>

-   <details>
    <summary>Tags API</summary>

    - **Public API**
        - List     
            - all
            - by topic id ( tags related to blogs under a specific topic )
    - **Private API**
        - Create
        - Update
        - Delete

    </details>

-   <details>
    <summary>Topics API</summary>

    - **Public API**
        - List     
    - **Private API**
        - Create
        - Update
        - Delete

    </details>

-   <details>
    <summary>Auth API</summary>

    - **Public API**
        - Login ( returns **JWT** token on success )
        - Logout ( removes **JWT** token from database )
        - Auth check ( mostly unused, checks if jwt token is valid )

    </details>

#### Details
> [Swaggo](https://github.com/swaggo/swag) (auto genrate swagger.json) dosn't support JWT auth yet. 
> The 'Authorization' header is placed as a headers param.
- Swagger Doc: [swagger.json](./backend/docs/swagger.json)

### More on backend
For more complete documantation on the backend, please refer to the [backend](./backend/) directory.

## CLI Tools
> Source code for cli tools are included in [backend](./backend) directory

### SyncTool
#### Installation
```bash
some command
```

I want to use my own editor to write notes.

This is a tool that can sync my notes to the server.

The notes should be organized like `dummyData` folder
- A **meta.yaml** containing tags and topics
- **blogs** folder containing blogs with frontmatter

After the first sync, an **ids.json** file will be created, which maps blog filenames to their ids.
This prevents blog ids from changing if we lost the database and need to sync from scratch.

#### About referential integrity
If some blog's frontmetter has a tag or topic that doesn't exist,
it will be recorded, logged out and written to a file (**data-inconsistency.json**), after which the sync process will be terminated.
Only after passing the validation process will the data be synced to the server.

### User register
#### Installation 
```bash
some command
```
> **This is build and placed alongside server binary in the docker image**

This project is only used by one person, with no intention of saving other user's stuff.

And because of this, there is no api for registoring a new user.

Only someone with direct access to the database can register.

#### Functions
- CRUD for user table, directly operates on the database.

## Frontend
Only used to display content.

Inspired by **Stack Overflow's** url design, slug is shown at the end of the url.

Ex: 
- `https://<domain>/blogs/<id>/<slug>`
- `https://notes.alexfangsw.com/blogs/1/a-dummy-blog-post`

### Tech stack
- Language: **Javascript**
- Framework: **[NextJS](https://nextjs.org/)**
- UI Libraries: 
    - **[Tailwindcss](https://tailwindcss.com/)**
    - **[Daisyui](https://daisyui.com/)**
