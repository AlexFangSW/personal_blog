import { merianda } from "@/app/fonts";
import Link from "next/link";
import Image from "next/image";
import pinIcon from "../../../../public/pin.svg";
import { LinkCard } from "@/app/components/linkCard";

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
async function Tags({ selected, topicID }) {
  // selected tag id current topic
  console.log(selected);
  const tags = [];

  const url = `${process.env.BACKEND_BASE_URL}/tags?topic=${topicID}`
  const res = await fetch(url)
  const parsedRes = await res.json()

  tags.push(
    <Link className="btn btn-ghost" href={`/topics/${topicID}`}>
      All
    </Link>,
  );

  for (const tag of parsedRes.msg) {
    tags.push(
      <Link className="btn btn-ghost" href={`/topics/${topicID}?tag=${tag.id}`}>
        {displayTagName(tag, selected)}
      </Link>,
    );
  }
  return tags;
}

function BlogTags({ tags }) {
  // list of tags
  const tag_list = [];
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">{tag.name}</div>);
  }
  return tag_list;
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
  const blogs = [];

  const url = new URL(`${process.env.BACKEND_BASE_URL}/blogs?topic=${topicID}`)
  if (tagID) {
    url.searchParams.append("tag", tagID)
  }
  const res = await fetch(url)
  const parsedRes = await res.json()


  for (const blog of parsedRes.msg) {
    blogs.push(
      <LinkCard href={`/blogs/${blog.id}`}>
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
    );
  }
  return blogs;
}

/**
  * Retruns info for the current topic
  *
  * @param {number} id 
  */
async function getCurrentTopic(id) {
  const url = `${process.env.BACKEND_BASE_URL}/topics/${id}`
  const res = await fetch(url)
  const parsedRes = await res.json()

  return parsedRes.msg
}

export default async function Page({ params, searchParams }) {
  console.log(searchParams["tag"]);
  const selectedTag = searchParams["tag"];
  const currentTopic = await getCurrentTopic(params.id)

  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          {currentTopic.name}
        </h1>
        <p>{currentTopic.description}</p>
        <div className="divider">Tags</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Tags selected={selectedTag} topicID={params.id} />
        </div>
        <div className="divider">Posts</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Blogs topicID={currentTopic.id} tagID={selectedTag} />
        </div>
      </div>
    </div>
  );
}
