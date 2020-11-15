const USER_STORAGE_KEY = "slack-user";

export default class UserStore {

  static saveUser(email) {
    window.localStorage.setItem(USER_STORAGE_KEY, email)
  }

  static getUser() {
    return  window.localStorage.getItem(USER_STORAGE_KEY)
  }

  static isUserSet() {
    const user = this.getUser();
    return user != null && typeof user !== 'undefined'
  }
}
