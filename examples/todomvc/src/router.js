const routes = [
  { path: '/', component: DashboardPage }
]

const router = new VueRouter({
  mode: 'history',
  routes // short for `routes: routes`
})
