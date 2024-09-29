import Link from "next/link";
import { merianda } from "./fonts";

export default function NotFound() {
  return (
    <div className="grid grid-cols-1 items-center justify-items-center min-h-screen p-5">
      <div className="flex flex-col items-center gap-y-5">
        <h1 className={`title text-6xl ${merianda.className}`}>
          Uncharted territory
        </h1>
        <div className="divider">404</div>
        <div className="flex items-center gap-x-5 ">
          <Link className="btn btn-neutral" href="/">
            Back to home
          </Link>
        </div>
      </div>
    </div>
  );
}
