import Vue from 'vue'
import Meta from 'vue-meta'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import 'material-design-icons-iconfont/dist/material-design-icons.css' // Ensure you are using css-loader
import { store } from './store/store'
import config from './config'

console.log({ RELAY_ADDRESS: config.relayAddress }) //eslint-disable-line

Vue.use(Meta)

Vue.config.productionTip = false

new Vue({
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
