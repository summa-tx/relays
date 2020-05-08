import Vue from 'vue'
import Vuex from 'vuex'
import external from './external'
import relay from './relay'

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules: {
    external,
    relay
  },

  state: {
    // For now, bitcoin net is always mainnet
    blockchainURL: 'https://blockstream.info/api'
  }
})
