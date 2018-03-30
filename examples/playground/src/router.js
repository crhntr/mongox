const routes = [
  { path: '/', component: DashboardPage },
  { path: '/c/:name', component: CollectionPage, props: true }
]

const router = new VueRouter({
  mode: 'history',
  routes // short for `routes: routes`
})
