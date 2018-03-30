const store = new Vuex.Store({
  plugins: [MonGoX.vuexPlugin()],
  state: {
    count: 0
  },
  getters: {
    count: state => {
      return state.count
    }
  },
  mutations: {
    increment (state) {
      state.count++
    }
  }
})
