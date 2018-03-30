var TodoInput = Vue.component('todo-input', {

  data: function () {
    return {
      description: ''
    }
  },

  template: `
  <input class="new-todo"
    autofocus autocomplete="off"
    placeholder="What needs to be done?"
    v-model="description"
    @keyup.enter="$store.dispatch('addTodo', {description})"/>`
})
