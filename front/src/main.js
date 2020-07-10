import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './plugins/element.js'
import './assets/iconfont/iconfont.css'
// 导入全局的样式表
import './assets/css/globle.css'
// import 'lib-flexible'
import 'lib-flexible/flexible'
import axios from 'axios'
import qs from 'qs'
// import http from './utils/http.js'
Vue.prototype.$qs = qs
// 配置请求的根路径
axios.defaults.baseURL = 'http://106.75.51.138:9001/api/v1/'
// axios.defaults.baseURL = 'http://localhost:9001/api/v1/'
// 为axios设置请求拦截器
Vue.prototype.$http = axios
// Vue.prototype.axios = http
// Vue.prototype.$http = http
Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
