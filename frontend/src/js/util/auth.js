export const getAuthUrlParams = () => {
  return window.location.search
}

export const getHappinessDoorId = () => {
  const params = extractQueryParams()
  return params.get('i')
}

const extractQueryParams = () => {
  return (new URL(document.location)).searchParams;
}