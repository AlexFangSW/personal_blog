import Link from "next/link"
import { merianda } from "../fonts"
function BlogNav() {
  return (
    <div className="navbar bg-neutral lg:items-start flex flex-col items-center  max-h-fit text-neutral-content">
      <Link className={`btn btn-ghost text-xl  ${merianda.className}`} href="/">
        CodingNotes
      </Link>
    </div>
  )
}

export { BlogNav }
