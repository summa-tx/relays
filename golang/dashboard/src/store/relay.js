import axios from 'axios'
import * as types from '@/store/mutation-types'
import { reverseEndianness } from '@/utils/utils'
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
  getBKD ({ commit, dispatch }) {
    axios.get(`${relayURL}/getbestdigest`)
      .then((res) => {
        commit(types.SET_CONNECTED, true)
        console.log('get BKD', res.data.result.result)
        const hashBE = reverseEndianness(res.data.result.result)
        console.log({hashBE})
        dispatch('info/setBKD', { hash: hashBE }, { root: true })
        dispatch('verifyHash', hashBE)
      })
      .catch((e) => {
        console.error('relay/getBKD:\n', e)
        if (e.message === 'Request failed with status code 500') {
          commit(types.SET_CONNECTED, false)
        }
      })
  },

  getLCA ({ commit, dispatch }) {
    axios.get(`${relayURL}/getlastreorglca`)
      .then((res) => {
        commit(types.SET_CONNECTED, true)
        console.log('get LCA', res.data.result.result)

        dispatch(
          'info/setLCA',
          {
            height: res.data.height,
            hash: reverseEndianness(res.data.result.result),
            verifiedAt: new Date()
          },
          { root: true }
        )
      })
      .catch((e) => {
        console.error('relay/getLCA:\n', e)
        commit(types.SET_CONNECTED, false)
      })
  },

  verifyHash ({ rootState, dispatch }, data) {
    console.log({ data })
    axios.get(`${rootState.blockchainURL}/blocks/${data.hashFromRelay}`)
      .then((block) => {
        console.log('block', block)
        dispatch(
          'info/setBKD', {
          height: block.data.height,
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
