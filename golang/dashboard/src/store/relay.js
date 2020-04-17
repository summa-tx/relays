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
      // Data structure:
      // {
      //   "height": "0",
      //   "result": {
      //     "result": "0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"
      //   }
      // }
      commit(types.GET_BKD, res.result.result)
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    })
  },

  getLCA ({ dispatch }) {
    axios.get(`${relayURL}/getlastreorglca`).then((res) => {
      // Data structure:
      // {
      //   "height": "0",
      //   "result": {
      //     "result": "0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"
      //   }
      // }
      dispatch('info/setRelayInfo', { key: 'lca', data: res })
      dispatch('info/setLastComms', { source: 'relay', date: new Date() })
    })
  },

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
