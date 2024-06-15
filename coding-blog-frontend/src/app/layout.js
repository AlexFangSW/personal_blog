import "./globals.css";
import { BlogFooter } from "./components/footer";
import { BlogNav } from "./components/nav";

export const metadata = {
  title: "Coding Notes",
  description: "A place to document things I have learned",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <BlogNav />
        {children}
        <BlogFooter />
      </body>
    </html>
  );
}
