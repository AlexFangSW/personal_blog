import Link from "next/link"
import { merianda } from "./fonts"
import { LinkCard } from "./components/linkCard"

async function Topics() {
  const topics = []

  const url = `${process.env.BACKEND_BASE_URL}/topics`
  const res = await fetch(url)
  const parsedRes = await res.json()

  if (parsedRes.status >= 500) {
    throw new Error(`home page load topics error: ${parsedRes.error}`)
  }

  for (const topic of parsedRes.msg) {
    topics.push(
      <LinkCard href={`/topics/${topic.id}/${topic.slug}`}>
        <div className="card-body">
          <h2 className="card-title">{topic.name}</h2>
          <p>{topic.description}</p>
        </div>
      </LinkCard>,
    )
  }
  return topics
}

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen items-center p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>Coding Notes</h1>
        <div className="flex items-center gap-x-5 ">
          <p>By: AlexFangSW</p>
          <div className="divider divider-horizontal"></div>
          <p>Email: alexfangsw@gmail.com</p>
          <Link
            className="btn btn-neutral"
            href="https://github.com/AlexFangSW"
          >
            GitHub
          </Link>
        </div>
      </div>
      <div className="divider">Topics</div>
      <div className="flex flex-row flex-wrap justify-center gap-2">
        <Topics />
      </div>
    </div>
  )
}
