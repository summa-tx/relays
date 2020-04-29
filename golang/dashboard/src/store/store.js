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
    // blockchainURL: process.env.MAINNET
    //   ? 'https://api.blockcypher.com/v1/btc/main'
    //   : 'https://api.blockcypher.com/v1/btc/test3',
    // For now, bitcoin net is always mainnet
    blockchainURL: 'https://api.blockcypher.com/v1/btc/main',
    extension: ''
  }
})
