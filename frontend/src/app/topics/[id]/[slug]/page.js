import { merianda } from "@/app/fonts"
import Link from "next/link"
import Image from "next/image"
import pinIcon from "../../../../../public/pin.svg"
import { LinkCard } from "@/app/components/linkCard"
import { getCurrentTopic } from "@/app/util/topic"
import { notFound, redirect } from 'next/navigation'

export const dynamic = 'force-dynamic'

/**
 * @param {object} tag
 * @param {number} selected 
 */
function displayTagName(tag, selected) {
  if (selected == tag.id) {
    return `[ ${tag.name} ]`
  }
  return tag.name
}

/**
  * Returns a list of tags relevent to the current topicID.
  * Selected topic will be marked.
  *
  * @param {Object} param0 
  * @param {number} param0.selected 
  * @param {number} param0.topicID 
  */
async function Tags({ selected, topic }) {
  // selected tag id current topic
  const tags = []

  const url = `${process.env.BACKEND_BASE_URL}/tags?topic=${topic.id}`
  const res = await fetch(url)
  const parsedRes = await res.json()

  if (parsedRes.status >= 500) {
    throw new Error(`topic page load tag error: ${parsedRes.error}`)
  }

  tags.push(
    <Link className="btn btn-ghost" href={`/topics/${topic.id}/${topic.slug}`}>
      All
    </Link>,
  )

  for (const tag of parsedRes.msg) {
    tags.push(
      <Link className="btn btn-ghost" href={`/topics/${topic.id}/${topic.slug}?tag=${tag.id}`}>
        {displayTagName(tag, selected)}
      </Link>,
    )
  }
  return tags
}

function BlogTags({ tags }) {
  // list of tags
  const tag_list = []
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">{tag.name}</div>)
  }
  return tag_list
}


/**
  * Returns a list of blog 'cards'.
  * Filtered by topic and tags.
  * @param {Object} param0 
  * @param {number} param0.topicID
  * @param {number} param0.tagID 
  */
async function Blogs({ topicID, tagID }) {
  // use single topic, multiple tags to filter blogs
  const blogs = []

  const url = new URL(`${process.env.BACKEND_BASE_URL}/blogs?topic=${topicID}`)
  if (tagID) {
    url.searchParams.append("tag", tagID)
  }
  const res = await fetch(url)
  const parsedRes = await res.json()

  if (parsedRes.status >= 500) {
    throw new Error(`topic page load blogs error: ${parsedRes.error}`)
  }

  for (const blog of parsedRes.msg) {
    blogs.push(
      <LinkCard href={`/blogs/${blog.id}/${blog.slug}`}>
        <div className="card-body">
          <div className="flex flex-row flex-wrap justify-left gap-2">
            {
              blog.pined ? (
                <Image priority className="h-5 w-5" src={pinIcon} alt="Pined" />
              ) : null
            }
            <div className="divider divider-horizontal"></div>
            <h2 className="card-title">{blog.title}</h2>
          </div>
          <p>{blog.description}</p>
          <div className="divider"></div>
          <div className="flex flex-row flex-wrap justify-left gap-2">
            <BlogTags tags={blog.tags} />
          </div>
        </div>
      </LinkCard>,
    )
  }
  return blogs
}


export default async function Page({ params, searchParams }) {
  const selectedTag = searchParams["tag"]
  const topicRes = await getCurrentTopic(params.id)

  if (topicRes.status >= 400) {
    notFound()
  } else if (topicRes.status >= 500) {
    throw new Error(`Topic page error: ${topicRes.error}`)
  }

  const currentTopic = topicRes.msg

  // adjust slug
  if (currentTopic.slug != params.slug) {
    redirect(`/topics/${params.id}/${currentTopic.slug}`)
  }

  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <h1 className={`title text-5xl ${merianda.className}`}>
        {currentTopic.name}
      </h1>
      <div className="flex flex-col w-full items-center gap-y-1 pt-5">
        <p>{currentTopic.description}</p>
        <div className="divider">Tags</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Tags selected={selectedTag} topic={currentTopic} />
        </div>
        <div className="divider">Posts</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Blogs topicID={currentTopic.id} tagID={selectedTag} />
        </div>
      </div>
    </div>
  )
}
