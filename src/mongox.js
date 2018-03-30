function createCollectionModule(name) {
  return {
    namespaced: true,

    actions: {
      insert: function (context, doc) {
        if (doc === undefined) {throw new Error('document undefined')}
        console.log(doc)
        return Vue.http.post(`/api/db/${name}/insert`, doc)
      } /* ,
      find: function (context, {query = {}, projection}) {
        var req = {query}
        if (projection !== undefined) {
          req.projection = projection
        }

        return Vue.http.post(`/api/db/${name}/find`, req)
      },
      update: function (context, {query, update}) {
        if (query === undefined) {throw new Error('query undefined')}
        if (update === undefined) {throw new Error('update undefined')}

        return Vue.http.post(`/api/db/${name}/update`, {query, update})
      },
      remove: function (context, {query}) {
        if (query === undefined) {throw new Error('query undefined')}

        return Vue.http.post(`/api/db/${name}/remove`, {query})
      } */
    },
    mutations: {}
  }
}


var dbModule = {
  namespaced: true,

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
        var cols = res.body.collections
        for (var i in cols) {
          store.registerModule(['db', cols[i]], createCollectionModule(cols[i]))
        }
        store.commit('db/collections', res.body.collections)
        console.log(store)
      }, err => {
        console.log(res.body)
      })
    }
  }
}
