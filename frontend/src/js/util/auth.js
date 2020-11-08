export const getAuthUrlParams = () => {
  return window.location.search
}

export const getHappinessDoorId = () => {
  const params = extractQueryParams()
  return params.get('i')
}

export const getJsonCookie = (name) => {
  const str = getCookieStr(name)
  if (str) {
    return JSON.parse(str)
  }
  return null
}

const getCookieStr = (name) => {
  const nameEQ = name + "=";
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

const extractQueryParams = () => {
  return (new URL(document.location)).searchParams;
}