import Vue from 'vue'
import Vuex from 'vuex'
import info from './info'

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules: {
    info
  },

  state: {
      isConnected: false,
      socketMessage: '',
      extension: ''
  },

  // TODO: Move all this to relay.js and convert to REST routes instead of socket calls. Update actions where used. Format of info dispatch actions should stay the same. See info.js for how data should be formatted. Will need to format BE and LE hex strings.

  // NB: All socket events follow the convention of SOCKET_<EVENT NAME>
  // The entire name must be capitalized
  // At this moment, I cannot figure out how to move sockets into its own
  // namespaced module. See main.js for instantiation. I've tried doing
  // 'sockets/SOCKET_' but that didn't work
  mutations: {
    RELAY_SOCKET_CONNECTED (state, connected) {
      state.isConnected = connected
    }
  },

  actions: {
    relay_socket_connect ({ commit, dispatch }) {
      console.log('connected to socket')
      commit('RELAY_SOCKET_CONNECTED', true)
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    },

    relay_socket_disconnect ({ commit }) {
      console.log('disconnected from socket')
      commit('RELAY_SOCKET_CONNECTED', false)
    },

    relay_socket_return_bkd ({ dispatch }, data) {
      console.log('return bkd')
      dispatch('info/setRelayInfo', { key: 'bkd', data })
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    },

    relay_socket_return_lca ({ dispatch }, data) {
      console.log('return lca')
      dispatch('info/setRelayInfo', { key: 'lca', data: data })
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    },

    relay_socket_new_extension({ state }, data) {
      console.log('new extension event', state.extension, data)
    },

    // Called when external info updates
    relay_socket_return_verify_height ({ dispatch }, data) {
      console.log('return verify height', data)
      if (data) {
        dispatch('info/setCurrentBlock', { verifiedAt: new Date() })
        dispatch('info/setLastComms', { source: 'relay', date: new Date() })
      }
    }
  }
})
