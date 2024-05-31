--name: CreateBlog: one

--name: ListBlogs: many

--name: ListBlogsDeleted: many

--name: ListBlogsAll: many

--name: GetBlog: one

--name: GetBlogDeleted: one

--name: UpdateBlog: one

--name: DeleteBlogSoft: one

--name: DeleteBlog: one

--name: RestoreBlog: one

-- hmmm... should I just return id or everything
--  - content might be very large
--  - will I actually use this retruned data ??
--  - I think I should create different functions, 
--    one returns no content, and another with content.

