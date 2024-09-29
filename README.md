# Coding Notes
[Build status XXX] 

Link to the website: [PLACE THE LINK HERE]

## Motive
A opportunity to practice **GO** while creating something I could use.

This project started around when **GO 1.22**  was released.
At that time I was choosing which framework or libraries I sould use. 
In the release notes, one thing cought my eye: 
**[Enhanced routing patterns](https://tip.golang.org/doc/go1.22#enhanced_routing_patterns)**.
This improved `net/http.ServeMux` quite a lot, adding features such as **path params**, **restrict methods to routes**...etc.

After seeing this plus a video on YouTube by **Dreams of Code**: [The standard library now has all you need for advanced routing in Go](https://www.youtube.com/watch?v=H7tbjKFSg58),
I thought, 'why not give the standard library a try ?' After than, this project was created.

## How it is used
I enjoy using my own editor, and I also want to seperate my content from the source code.

- Deploy the server, and create a user with **[UserRegister](./backend/README.md#userregister)** CLI tool.
- Write blogs and organize it like [dummyData](./backend/dummyData/) directory and store is in a seperated repository.
- When pushed, it triggers a CI/CD pipeline that uses **[SyncTool](./backend/README.md#synctool)** to sync data to the server.


## Project Stucture
This project is separated into **frontend**, **backend** and **CLI tools** for me to sync my data to the server.
```
.
├── backend <--- Backend and CLI tools are both in this folder
├── frontend
└── README.md
```

## Backend
This part is mostly made with the **GO** standard library,
not just routing, but also database related stuff.

### Tech stack
- Language: **GO**
- Server: **net/http** ( Routing is done with `net/http.ServeMux`)
- ORM: **None**, this project uses `database/sql` with **raw SQL queries**.
- Database: **SQLite**
- Database Migration: **[goose](https://github.com/pressly/goose)**
- API documentation: generated with **[swaggo](https://github.com/swaggo/swag)**

For more documentation on the backend, please refer to the [backend](./backend/) directory.

## Frontend
Inspired by **Stack Overflow's** URL design, slug is shown at the end of the url.

Ex: 
- `https://<domain>/blogs/<id>/<slug>`
- `https://notes.alexfangsw.com/blogs/1/a-dummy-blog-post`

### Tech stack
- Language: **Javascript**
- Framework: **[NextJS](https://nextjs.org/)**
- UI Libraries: 
    - **[Tailwindcss](https://tailwindcss.com/)**
    - **[Daisyui](https://daisyui.com/)**

## CLI Tools
- UserRegister
- SyncTool

Documentation are at [CLI tools](./backend/README.md#cli-tools) section.
