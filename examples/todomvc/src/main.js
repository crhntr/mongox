function main () {z
  new Vue({
    el: '#app',
    data: function () {
      return {
        msg: 'hello'
      }
    },
    router,
    store,
    directives: {
      'todo-focus': function (el, binding) {
        if (binding.value) {
          el.focus()
        }
      }
    }
  })
}
