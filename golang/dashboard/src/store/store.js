import Vue from 'vue'
import Vuex from 'vuex'
import info from './info'
import relay from './relay'

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules: {
    info,
    relay
  },

  state: {
    extension: ''
  }
})
