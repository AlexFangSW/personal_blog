import Link from "next/link"

function LinkCard({ href, children }) {
  return (
    <Link href={href} className="card lg:w-96 w-11/12 bg-neutral text-neutral-content">
      {children}
    </Link>
  )
}

export { LinkCard }
