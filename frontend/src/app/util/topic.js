/**
  * Retruns info for the current topic
  *
  * @param {number} id 
  */
async function getCurrentTopic(id) {
  const url = `${process.env.BACKEND_BASE_URL}/topics/${id}`
  const res = await fetch(url)
  const parsedRes = await res.json()

  return parsedRes
}

export { getCurrentTopic }
