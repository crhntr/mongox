var TodoInput = Vue.component('todo-input', {

  data: function () {
    return {
      description: ''
    }
  },

  methods: {
    addTodo: function () {
      console.log(this.description)
    }
  },

  template: `
  <input class="new-todo"
    autofocus autocomplete="off"
    placeholder="What needs to be done?"
    v-model="description"
    @keyup.enter="addTodo"/>`
})
