import axios from 'axios'
// import * as types from '@/store/mutation-types'
const relayURL = '/relay'

// const state = {

// }

// const mutations: {

// }

// TODO: Convert to REST routes instead of socket calls. Update actions where used. Format of info dispatch actions should stay the same. See info.js for how data should be formatted. Will need to format BE and LE hex strings.

const actions = {
  getBKD ({ dispatch }) {
    axios.get(`${relayURL}/getbestdigest`).then((res) => {
      console.log('get BKD', res)
      // Data structure:
      // {
      //   "height": "0",
      //   "result": {
      //     "result": "0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"
      //   }
      // }
      // commit(types.GET_BKD, res.data.result.result)
      dispatch(
        'info/setRelayInfo',
        {
          key: 'bkd',
          // TODO: switch endian
          data: res.data.result.result
        },
        { root: true }
      )
      dispatch(
        'info/setLastComms',
        { source: 'relay', date: new Date() },
        { root: true }
      )
    })
    .catch((e) => {
      console.log('relay/getBKD: ', e)
    })
  },

  getLCA ({ dispatch }) {
    axios.get(`${relayURL}/getlastreorglca`).then((res) => {
      console.log('get LCA', res)
      // Data structure:
      // {
      //   "height": "0",
      //   "result": {
      //     "result": "0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"
      //   }
      // }
      dispatch(
        'info/setRelayInfo',
        {
          key: 'lca',
          // TODO: switch endian
          data: res.data.result.result
        },
        { root: true }
      )
      dispatch(
        'info/setLastComms',
        { source: 'relay', date: new Date() },
        { root: true }
      )
    })
  },

  // relay_socket_new_extension({ state }, data) {
  //   console.log('new extension event', state.extension, data)
  // },

  // NB: Verify height does not actually verify height. This is for updating only. See `info/getExternalInfo`
  verifyHeight ({ dispatch }, data) {
    console.log('verify height', data)
    if (data) {
      dispatch(
        'info/setCurrentBlock',
        { verifiedAt: new Date() },
        { root: true }
      )
      dispatch(
        'info/setLastComms',
        { source: 'relay', date: new Date() },
        { root: true }
      )
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
