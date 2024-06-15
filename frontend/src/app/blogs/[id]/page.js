import { merianda } from "@/app/fonts";
import Link from "next/link";
import { MDXRemote } from "next-mdx-remote/rsc";
import { promises as fs } from "fs";
function BlogTags({ tags }) {
  // list of tags
  const tag_list = [];
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">{tag}</div>);
  }
  return tag_list;
}

function BlogTopics({ topics }) {
  // list of topics
  const topic_list = [];
  for (const topic of topics) {
    topic_list.push(
      <Link className="badge badge-outline" href={`/topics/${topic}`}>
        {topic}
      </Link>,
    );
  }
  return topic_list;
}

async function getBlogData(id) {
  // MDX text - can be from a local file, database, CMS, fetch, anywhere...
  // const res = await fetch('https://...')
  // const markdown = await res.text()

  // NOTE: markdown rander can't handle '{}' in content, unless we escape them with '\'
  const file = await fs.readFile(
    process.cwd() + "/src/app/test_blog.md",
    "utf8",
  );
  const data = file;
  return { data };
}

export default async function Page({ params }) {
  const { data } = await getBlogData(params.id);

  const customizeComponents = {
    p: (props) => <p className="text-xl">{props.children}</p>,
  };

  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          Blog number: {params.id}
        </h1>
        <div className="flex flex-row flex-wrap justify-center items-center gap-2">
          <p>By: AlexFangSW</p>
          <div className="divider divider-horizontal"></div>
          <div className="flex flex-col flex-wrap justify-center gap-2">
            <p>Created at: 2024-06-14T11:39:56+00:00</p>
            <p>Updated at: 2024-06-14T11:39:56+00:00</p>
          </div>
        </div>
        <p>blog description aaaa bbbb cccc fjdskalfjdkslajfkdlsa</p>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <BlogTopics topics={["topic12"]} />
          <div className="divider divider-horizontal"></div>
          <BlogTags tags={["aaa", "aaa", "aaa"]} />
        </div>
        <div className="divider">Content</div>
        <article className="prose max-w-4xl  text-neutral-content">
          <MDXRemote source={data} components={customizeComponents} />
        </article>
      </div>
    </div>
  );
}
