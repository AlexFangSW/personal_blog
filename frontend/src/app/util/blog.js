/**
 * @param {int} id 
 */
async function getCurrentBlog(id) {
  const url = `${process.env.BACKEND_BASE_URL}/blogs/${id}?parsed=true`
  const res = await fetch(url)
  const parsedRes = await res.json()

  return parsedRes
}

export { getCurrentBlog }
