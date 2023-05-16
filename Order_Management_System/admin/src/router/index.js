import Vue from 'vue'
import VueRouter from 'vue-router'
import HomeView from '../views/HomeView.vue'
import MainView from '../views/MainView.vue'
import AboutView from '../views/AboutView.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'MainView',
    component: MainView,
    children: [
      {path: '/home/about', component: AboutView}
    ]
  },
  {
    path: '/home',
    name: 'HomeView' ,
    component: HomeView
  }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
