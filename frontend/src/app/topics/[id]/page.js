import { merianda } from "@/app/fonts";
import Link from "next/link";
import Image from "next/image";
import pinIcon from "../../../../public/pin.svg";

async function Tags({ selected, topicID }) {
  // selected tag id current topic
  console.log(selected);
  const tags = [];
  tags.push(
    <Link className="btn btn-ghost" href={`/topics/${topicID}`}>
      All
    </Link>,
  );
  for (let index = 0; index < 10; index++) {
    tags.push(
      <Link className="btn btn-ghost" href={`/topics/${topicID}?tag=${index}`}>
        Ghost {selected}
      </Link>,
    );
  }
  return tags;
}

function BlogTags({ tags }) {
  // list of tags
  const tag_list = [];
  for (const tag of tags) {
    tag_list.push(<div className="badge badge-outline">{tag}</div>);
  }
  return tag_list;
}

async function Blogs({ topic, tags }) {
  // use single topic, multiple tags to filter blogs
  const blogs = [];
  for (let index = 0; index < 10; index++) {
    blogs.push(
      <Link
        className="card w-96 bg-neutral text-neutral-content"
        href={`/blogs/${index}`}
      >
        <div className="card-body">
          <div className="flex flex-row flex-wrap justify-left gap-2">
            <Image priority className="h-5 w-5" src={pinIcon} alt="Pined" />
            <div className="divider divider-horizontal"></div>
            <h2 className="card-title">Blog Number: {index}</h2>
          </div>
          <p>
            If a dog chews shoes whose shoes does he choose?
            fjdkslajfkdasjfkdlsajfldsajfljlsajfsak
          </p>
          <div className="divider"></div>
          <div className="flex flex-row flex-wrap justify-left gap-2">
            <BlogTags tags={["aaa", "bbb", "aaa", "bbb"]} />
          </div>
        </div>
      </Link>,
    );
  }
  return blogs;
}

export default function Page({ params, searchParams }) {
  console.log(searchParams["tag"]);
  const selectedTag = searchParams["tag"];
  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          Topic number: {params.id}
        </h1>
        <p>topics description aaaa bbbb cccc fjdskalfjdkslajfkdlsa</p>
        <div className="divider">Tags</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Tags selected={selectedTag} topicID={params.id} />
        </div>
        <div className="divider">Posts</div>
        <div className="flex flex-row flex-wrap justify-center gap-2">
          <Blogs />
        </div>
      </div>
    </div>
  );
}
