import Image from "next/image"
import githubIcon from "../../../public/github.svg"
import Link from "next/link"
function BlogFooter() {
  const currentYear = new Date().getFullYear()
  return (
    <footer className="footer items-center gap-0 max-h-fit p-1 pl-5 bg-neutral text-neutral-content">
      <aside className="items-center lg:grid-flow-col lg:gap-2 gap-0">
        <p>Copyright © {currentYear} - All right reserved</p>
        <p>By: AlexFangSW</p>
      </aside>
      <nav className="grid-flow-col lg:gap-4 lg:place-self-center lg:justify-self-end items-center">
        <p>Email: alexfangsw@gmail.com</p>
        <Link href="https://github.com/AlexFangSW" >
          <Image priority className="w-8 hidden lg:block" src={githubIcon} alt="GitHub" />
        </Link>
      </nav>
    </footer>
  )
}

export { BlogFooter }
