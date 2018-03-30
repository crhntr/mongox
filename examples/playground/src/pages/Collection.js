var CollectionPage = Vue.component('collection-page', {
  props: ['name'],
  template: `<section>
  <h1>{{name}}</h1>
  <create-document></create-document>
</section>`
})
