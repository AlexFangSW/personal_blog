import Link from "next/link"
import { merianda } from "../fonts"
function BlogNav() {
  return (
    <div className="navbar bg-neutral text-neutral-content">
      <Link className={`btn btn-ghost md:text-xl  ${merianda.className}`} href="/">
        CodingNotes
      </Link>
    </div>
  )
}

export { BlogNav }
