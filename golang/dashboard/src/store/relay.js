import axios from 'axios'
import * as types from '@/store/mutation-types'
import { remove0x, reverseEndianness } from '@/utils/utils'
const relayURL = '/relay'

const state = {
  connected: true
}

const mutations = {
  [types.SET_CONNECTED] (state, connected) {
    state.connected = connected
  }
}

const actions = {
  getBKD ({ rootState, commit, dispatch }) {
    axios.get(`${relayURL}/getbestdigest`).then((res) => {
      console.log('get BKD', res)
      commit(types.SET_CONNECTED, true)

      const hashBE = reverseEndianness(res.data.result.result)
      console.log('blockchainURL', rootState.blockchainURL)
      axios.get(`${rootState.blockchainURL}/blocks/${remove0x(hashBE)}`).then((block) => {
        console.log('block', block)
        dispatch(
          'info/setBKD',
          {
            height: block.data.height,
            hash: hashBE,
            verifiedAt: new Date()
          },
          { root: true }
        )
        dispatch(
          'info/setLastComms',
          { source: 'relay', date: new Date() },
          { root: true }
        )
      }).catch((e) => {
        console.error('relay/getBKD:\n', e)
      })

      // dispatch(
      //   'info/setBKD',
      //   {
      //     height: res.data.height,
      //     hash: hashBE,
      //     verifiedAt: new Date()
      //   },
      //   { root: true }
      // )
      // dispatch(
      //   'info/setLastComms',
      //   { source: 'relay', date: new Date() },
      //   { root: true }
      // )
    })
    .catch((e) => {
      console.error('relay/getBKD:\n', e)
      commit(types.SET_CONNECTED, false)
    })
  },

  getLCA ({ commit, dispatch }) {
    axios.get(`${relayURL}/getlastreorglca`).then((res) => {
      console.log('get LCA', res)
      commit(types.SET_CONNECTED, true)
      // Data structure:
      // {
      //   "height": "0",
      //   "result": {
      //     "result": "0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"
      //   }
      // }
      dispatch(
        'info/setLCA',
        {
          height: res.data.height,
          hash: reverseEndianness(res.data.result.result),
          verifiedAt: new Date()
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
      console.error('relay/getLCA:\n', e)
      commit(types.SET_CONNECTED, false)
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
}

// state,
// mutations,
export default {
  namespaced: true,
  state,
  mutations,
  actions
}
