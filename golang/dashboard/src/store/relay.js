import axios from 'axios'
import * as types from '@/store/mutation-types'
import {
  reverseEndianness,
  convertUnixTimestamp,
  lStorage
} from '@/utils/utils'
const relayURL = '/relay'

const state = {
  connected: true,

  lastComms: lStorage.get('lastCommsRelay') || undefined,

  // Best Known Digest
  bkd: lStorage.get('bkd') || {
    height: 0,              // Number - height of the BKD
    hash: '',               // String - BKD hash
    time: undefined,        // Date - BKD timestamp, from external
    updatedAt: undefined   // Date - When was the BKD last verified
  },

  // last (reorg) common ancestor
  lca: lStorage.get('lca') || {
    height: 0,              // Number - height of the LCA
    hash: '',               // String - LCA hash
    time: undefined,        // Date - LCA timestamp, from external
    updatedAt: undefined   // Date - When was the LCA last verified
  }
}

const mutations = {
  [types.SET_CONNECTED] (state, connected) {
    state.connected = connected
  },

  [types.SET_LAST_COMMS] (state, { date }) {
    state.lastComms = date
    lStorage.set('lastCommsRelay', state.lastComms)
  },

  // NB: BKD = best known digest
  [types.SET_BKD] (state, payload) {
    for (let key in payload) {
      state.bkd[key] = payload[key]
    }
    lStorage.set('bkd', state.bkd)
  },

  // NB: LCA = last (reorg) common ancestor
  [types.SET_LCA] (state, payload) {
    for (let key in payload) {
      state.lca[key] = payload[key]
    }
    lStorage.set('lca', state.lca)
  }
}

const actions = {
  getBKD ({ commit, dispatch }) {
    axios.get(`${relayURL}/getbestdigest`)
      .then((res) => {
        commit(types.SET_CONNECTED, true)

        const hashBE = reverseEndianness(res.data.result.result)
        console.log('get BKD: ', hashBE)

        dispatch('setBKD', { hash: hashBE })
        dispatch('verifyHash', { hash: hashBE, type: 'BKD' })
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

        const hashBE = reverseEndianness(res.data.result.result)
        console.log('get LCA: ', hashBE)

        dispatch('setLCA', { hash: hashBE })
        dispatch('verifyHash', { hash: hashBE, type: 'LCA'})
      })
      .catch((e) => {
        console.error('relay/getLCA:\n', e)
        if (e.message === 'Request failed with status code 500') {
          commit(types.SET_CONNECTED, false)
        }
      })
  },

  verifyHash ({ rootState, dispatch }, data) {
    // data.hash, data.type = 'BKD', 'LCA'
    console.log({ data })
    axios.get(`${rootState.blockchainURL}/block/${data.hash}`)
      .then((block) => {
        console.log('block', block)
        dispatch(
          `set${data.type}`,
          {
            height: block.data.height,
            time: convertUnixTimestamp(block.data.timestamp),
            updatedAt: new Date()
          }
        )
        dispatch(
          'info/setLastComms',
          { date: new Date() },
          { root: true }
        )
      }).catch((e) => {
        console.error('relay/verifyHash:\n', e)
      })
  },

  // payload: { key: '', data: '' }
  setBKD ({ commit }, payload) {
    commit(types.SET_BKD, payload)
    commit(types.SET_LAST_COMMS, { date: new Date() })
  },

  // payload: { key: '', data: '' }
  setLCA ({ commit }, payload) {
    commit(types.SET_LCA, payload)
    commit(types.SET_LAST_COMMS, { date: new Date() })
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
