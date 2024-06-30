import { getCurrentBlog } from '@/app/util/blog';
import { redirect } from 'next/navigation'
import { notFound } from 'next/navigation';

export default async function ToBlog({ params }) {
  const blogRes = await getCurrentBlog(params.id)

  if (blogRes.status == 404) {
    notFound()
  } else if (blogRes.status >= 500) {
    throw new Error(`To blog page error: ${blogRes.error}`)
  }

  const currentBlog = blogRes.msg
  redirect(`/blogs/${params.id}/${currentBlog.slug}`)
}
