'use client'

import Link from "next/link";
import { merianda } from "./fonts";
import { useEffect } from "react";

export default function ErrorPage({ error, reset }) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div className="grid grid-cols-1 items-center justify-items-center min-h-screen p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          Error!!!
        </h1>
        <div className="divider">500</div>
        <div className="flex items-center gap-x-5 ">
          <button className="btn btn-neutral"
            onClick={
              () => reset()
            }
          >
            Try again
          </button>
          <Link className="btn btn-neutral" href="/">
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
