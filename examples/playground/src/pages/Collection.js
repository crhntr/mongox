var CollectionPage = Vue.component('collection-page', {
  props: ['name'],
  template: `<section>
  <h1>{{name}}</h1>
  <insert-json :col="name"></insert-json>
</section>`
})
