export default function initMongox (url, {}) {

  return store => {

    store.subscribe((mutation, state) => {
      // called after every mutation.
      // The mutation comes in the format of `{ type, payload }`.
    })
  }
}
