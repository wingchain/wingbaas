import Vue from 'vue'
import Vuex from 'vuex'
// 持久化store里面的数据
import createPersistedState from 'vuex-persistedstate'
Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    title: '我的联盟',
    isShow: ''
  },
  mutations: {
    changeTitle (state, Title) {
      state.title = Title
      state.isShow = true
    },
    changeshow (state, Title) {
      state.title = Title
      state.isShow = false
    }
  },
  actions: {
  },
  modules: {
  },
  plugins: [createPersistedState({
    // storage: window.sessionStorage,
    // reducer(data) {
    //  return {
    //  // 设置只储存state中的myData
    //  myData: data.myData
    // }
    // }
  })]
})
