import "./globals.css";
import { BlogFooter } from "./components/footer";
import { BlogNav } from "./components/nav";
import { ErrorBoundary } from "next/dist/client/components/error-boundary";
import ErrorPage from "./error";

export const metadata = {
  title: "Coding Notes",
  description: "A place to document things I have learned",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <BlogNav />
        <ErrorBoundary fallback={<ErrorPage />} >
          {children}
        </ErrorBoundary>
        <BlogFooter />
      </body>
    </html>
  );
}
