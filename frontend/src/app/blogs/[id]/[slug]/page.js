import { merianda } from "@/app/fonts"
import Link from "next/link"
import { MDXRemote } from "next-mdx-remote/rsc"
import { getCurrentBlog } from "@/app/util/blog"
import { notFound, redirect } from "next/navigation"

/**
 * @param {Object} param0 
 * @param {Object} param0.tags 
 */
function BlogTags({ tags }) {
  const tag_list = []
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">
      {tag.name}
    </div>)
  }
  return tag_list
}

/**
 * @param {Object} param0 
 * @param {Object} param0.topics 
 */
function BlogTopics({ topics }) {
  const topic_list = []
  for (const topic of topics) {
    topic_list.push(
      <Link className="badge badge-outline" href={`/topics/${topic.id}/${topic.slug}`}>
        {topic.name}
      </Link>,
    )
  }
  return topic_list
}


export default async function Page({ params }) {
  const blogRes = await getCurrentBlog(params.id)

  if (blogRes.status >= 400) {
    notFound()
  } else if (blogRes.status >= 500) {
    throw new Error(`Blog page error: ${blogRes.error}`)
  }

  const currentBlog = blogRes.msg

  const createAtDate = new Date(currentBlog.created_at)
  const createAtStr = `${createAtDate.getFullYear()}/${createAtDate.getMonth()}/${createAtDate.getDay()}`
  const updateAt = new Date(currentBlog.updated_at)
  const updateAtStr = `${updateAt.getFullYear()}/${updateAt.getMonth()}/${updateAt.getDay()}`

  // adjust slug
  if (currentBlog.slug != params.slug) {
    redirect(`/blogs/${params.id}/${currentBlog.slug}`)
  }

  const customizeComponents = {
    p: (props) => <p className="text-xl">{props.children}</p>,
    strong: (props) => <strong className="text-inherit">{props.children}</strong>,
  }

  // NOTE: markdown rander can't handle '{}' in content, unless we escape them with '\'
  return (
    <div className="flex flex-col min-h-screen  items-center p-5">
      <div className="flex flex-col items-center gap-y-10">
        <h1 className={`title text-6xl ${merianda.className}`}>
          {currentBlog.title}
        </h1>
        <div className="flex flex-row flex-wrap justify-center items-center gap-2">
          <p>Created at: {createAtStr}</p>
          <div className="divider divider-horizontal"></div>
          <p>By: AlexFangSW</p>
          <div className="divider divider-horizontal"></div>
          <p>Updated at: {updateAtStr}</p>
        </div>
      </div>
      <div className="flex flex-col items-center gap-y-5 pt-3">
        <p>{currentBlog.description}</p>
        <div className="flex flex-row flex-wrap justify-center gap-3">
          <BlogTopics topics={currentBlog.topics} />
          <div className="divider divider-horizontal"></div>
          <BlogTags tags={currentBlog.tags} />
        </div>
        <div className="divider">Content</div>
        <article className="prose max-w-4xl text-neutral-content">
          <MDXRemote source={currentBlog.content} components={customizeComponents} />
        </article>
      </div>
    </div >
  )
}
