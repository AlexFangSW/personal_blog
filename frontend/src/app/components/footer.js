import Image from "next/image";
import githubIcon from "../../../public/github.svg";
import Link from "next/link";
function BlogFooter() {
  const currentYear = new Date().getFullYear();
  return (
    <footer className="footer items-center p-4 bg-neutral text-neutral-content">
      <aside className="items-center grid-flow-col">
        <p>Copyright © {currentYear} - All right reserved</p>
        <div className="divider divider-horizontal"></div>
        <p>By: AlexFangSW</p>
      </aside>
      <nav className="grid-flow-col gap-4 md:place-self-center md:justify-self-end items-center">
        <p>Email: alexfangsw@gmail.com</p>
        <Link href="https://github.com/AlexFangSW">
          <Image priority className="w-10" src={githubIcon} alt="GitHub" />
        </Link>
      </nav>
    </footer>
  );
}

export { BlogFooter };
