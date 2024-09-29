import { getCurrentTopic } from '@/app/util/topic'
import { redirect } from 'next/navigation'
import { notFound } from 'next/navigation';

export default async function ToTopic({ params }) {
  const topicRes = await getCurrentTopic(params.id)

  if (topicRes.status == 404) {
    notFound()
  } else if (topicRes.status >= 500) {
    throw new Error(`To topic page error: ${topicRes.error}`)
  }

  const currentTopic = topicRes.msg
  redirect(`/topics/${params.id}/${currentTopic.slug}`)
}
