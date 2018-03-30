function createCollectionModule(name) {
  return {
    namespaced: true,

    actions: {
      insert: function (doc) {
        if (doc === undefined) {throw new Error('document undefined')}

        return Vue.http.post(`/api/db/${name}/insert`, req)
      },
      find: function (query = {}, projection) {
        var req = {query}
        if (projection !== undefined) {
          req.projection = projection
        }

        return Vue.http.post(`/api/db/${name}/find`, req)
      },
      update: function (query, update) {
        if (query === undefined) {throw new Error('query undefined')}
        if (update === undefined) {throw new Error('update undefined')}

        return Vue.http.post(`/api/db/${name}/update`, {query, update})
      },
      remove: function (query) {
        if (query === undefined) {throw new Error('query undefined')}

        return Vue.http.post(`/api/db/${name}/remove`, {query})
      }
    },
    mutations: {}
  }
}


var dbModule = {
  store: {
    collections: []
  },

  mutations: {
    collections: (state, cols) => {
      state.collections = cols
    }
  },

  getters: {
    collections: state => {
      return state.collections
    }
  }
}

class MonGoX {
  static vuexPlugin () {
    return store => {
      store.registerModule('db', dbModule)
      Vue.http.get('/api/db').then(res => {
        var cols = res.body.collections || []
        for (var col in cols) {
          store.registerModule(['db', col], createCollectionModule(col))
        }
        store.commit('collections', cols)
        console.log(store)
      }, err => {
        console.log(res.body)
      })
    }
  }
}
