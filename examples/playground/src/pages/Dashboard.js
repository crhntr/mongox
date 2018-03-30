var DashboardPage = Vue.component('dashboard-page', {
  template: `<section>
  <nav>
    <router-link v-for="col in $store.getters.collections" :to="'/c/'+col" :key="col">{{col}}</router-link>
  </nav>
</section>`
})
