import { merianda } from "@/app/fonts";
import Link from "next/link";
import { MDXRemote } from "next-mdx-remote/rsc";

/**
 * @param {Object} param0 
 * @param {Object} param0.tags 
 */
function BlogTags({ tags }) {
  const tag_list = [];
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">
      {tag.name}
    </div>);
  }
  return tag_list;
}

/**
 * @param {Object} param0 
 * @param {Object} param0.topics 
 */
function BlogTopics({ topics }) {
  const topic_list = [];
  for (const topic of topics) {
    topic_list.push(
      <Link className="badge badge-outline" href={`/topics/${topic.id}/${topic.slug}`}>
        {topic.name}
      </Link>,
    );
  }
  return topic_list;
}

/**
 * @param {int} id 
 */
async function getBlogData(id) {
  const url = `${process.env.BACKEND_BASE_URL}/blogs/${id}`
  const res = await fetch(url)
  const parsedRes = await res.json()

  // NOTE: markdown rander can't handle '{}' in content, unless we escape them with '\'
  return parsedRes.msg;
}

export default async function Page({ params }) {
  const blog = await getBlogData(params.id);

  const customizeComponents = {
    p: (props) => <p className="text-xl">{props.children}</p>,
    strong: (props) => <strong className="text-inherit">{props.children}</strong>,
  };

  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          {blog.title}
        </h1>
        <div className="flex flex-row flex-wrap justify-center items-center gap-2">
          <p>By: AlexFangSW</p>
          <div className="divider divider-horizontal"></div>
          <div className="flex flex-col flex-wrap justify-center gap-2">
            <p>Created at: {blog.created_at}</p>
            <p>Updated at: {blog.updated_at}</p>
          </div>
        </div>
        <p>{blog.description}</p>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <BlogTopics topics={blog.topics} />
          <div className="divider divider-horizontal"></div>
          <BlogTags tags={blog.tags} />
        </div>
        <div className="divider">Content</div>
        <article className="prose max-w-4xl  text-neutral-content">
          <MDXRemote source={blog.content} components={customizeComponents} />
        </article>
      </div>
    </div>
  );
}
