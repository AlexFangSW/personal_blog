import Link from "next/link";

function LinkCard({ href, children }) {
  return (
    <Link href={href} className="card w-96 bg-neutral text-neutral-content">
      {children}
    </Link>
  );
}

export { LinkCard };
