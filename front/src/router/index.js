import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../components/login'
import Registered from '../components/registered'
import Home from '../components/home'
import Create from '../components/createaliance/create'
import Cluster from '../components/clusterManage/cluster'
import Deploy from '../components/createDeploy/deploy'
import Treaty from '../components/uploadocc/treaty'
import User from '../components/joinuser/user'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    component: Login
  },
  {
    path: '/registered',
    component: Registered
  },
  {
    path: '/home',
    component: Home,
    redirect: '/create',
    children: [{
      path: '/create',
      component: Create
    },
    {
      path: '/cluster',
      component: Cluster
    },
    {
      path: '/deploy',
      component: Deploy
    },
    {
      path: '/treaty',
      component: Treaty
    },
    {
      path: '/user',
      component: User
    }
    ]
  }
]
// 解决ElementUI导航栏中的vue-router在3.0版本以上重复点菜单报错问题
const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function push (location) {
  return originalPush.call(this, location).catch(err => err)
}
const router = new VueRouter({
  routes
})
// 利用导航守卫控制访问权限
// router.beforeEach((to, from, next) => {
//   if (to.path === '/login') return next()
//   const tokenStr = window.sessionStorage.getItem('token')
//   if (!tokenStr) return next('/login')
//   next()
// })

export default router
