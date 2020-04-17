import axios from 'axios'
import * as types from '@/store/mutation-types'
const relayURL = 'http://localhost:1317/relay'

// const state = {
  
// }

// const mutations: {

// }

// TODO: Convert to REST routes instead of socket calls. Update actions where used. Format of info dispatch actions should stay the same. See info.js for how data should be formatted. Will need to format BE and LE hex strings.

const actions = {
  getBKD ({ commit, dispatch }) {
    axios.get(`${relayURL}/getbestdigest`).then((res) => {
      commit(types.GET_BKD, res)
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    })
  },

  // relay_socket_return_bkd ({ dispatch }, data) {
    //   console.log('return bkd')
    //   dispatch('info/setRelayInfo', { key: 'bkd', data })
    //   dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    // },

  getLCA ({ dispatch }) {
    axios.get(`${relayURL}/getlastreorglca`).then((res) => {
      dispatch('info/setRelayInfo', { key: 'lca', data: res })
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    })
  },
  // relay_socket_return_lca ({ dispatch }, data) {
  //   console.log('return lca')
  //   dispatch('info/setRelayInfo', { key: 'lca', data: data })
  //   dispatch('info/setLastComms', { source: 'relay', date: new Date() })
  // },

  // relay_socket_new_extension({ state }, data) {
  //   console.log('new extension event', state.extension, data)
  // },

  verifyHeight ({ dispatch }, data) {
    if (data) {
      dispatch('info/setCurrentBlock', { verifiedAt: new Date() })
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    }
  }

  // // Called when external info updates
  // relay_socket_return_verify_height ({ dispatch }, data) {
  //   console.log('return verify height', data)
  //   if (data) {
  //     dispatch('info/setCurrentBlock', { verifiedAt: new Date() })
  //     dispatch('info/setLastComms', { source: 'relay', date: new Date() })
  //   }
  // }
}

// state,
// mutations,
export default {
  namespaced: true,
  actions
}
