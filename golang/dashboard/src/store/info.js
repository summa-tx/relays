import axios from 'axios'
import * as types from '@/store/mutation-types'
import { lStorage } from '@/utils/utils'
const isMain = process.env.MAINNET
const blockchainURL = isMain
  ? 'https://api.blockcypher.com/v1/btc/main'
  : 'https://api.blockcypher.com/v1/btc/test3'

const state = {
  source: 'blockcypher.com',
  // When was the last time a successful communication was made?
  // If last successful check was > 5 minutes ago, show flag
  lastComms: lStorage.get('lastComms') || {
    relay: undefined,       // Date
    external: undefined,    // Date
  },

  currentBlock: lStorage.get('currentBlock') || {
    height: 0,              // Number - Current block height, from external
    hash: '',               // String - Current block hash, from external
    updatedAt: undefined,   // Date - When was this data updated
    verifiedAt: undefined   // Date - When block was verified
  },

  // Keep track of previous block information
  // If incoming block number increments, then move currentBlock info to here
  // and incoming block info goes to currentBlock
  previousBlocks: lStorage.get('previousBlocks') || [],

  // Relay
  relay: lStorage.get('relay') || {
    bkd: '',     // String - best known digest
    lca: ''      // String - last (reorg) common ancestor
  }
}

const mutations = {
  [types.SET_LAST_COMMS] (state, { source, date }) {
    state.lastComms[source] = date
    lStorage.set('lastComms', state.lastComms)
  },

  [types.SET_CURRENT_BLOCK] (state, block) {
    let newBlock = state.currentBlock
    Object.keys(block).forEach((prop) => {
      newBlock[prop] = block[prop]
    })
    state.currentBlock = newBlock
    lStorage.set('currentBlock', state.currentBlock)
  },

  // This is called when current block is updated
  // Take all data and put it here
  // TODO: Make sure to control and handle duplicates
  [types.ADD_PREVIOUS_BLOCK] (state, block) {
    state.previousBlocks.push(block)
    lStorage.set('previousBlocks', state.previousBlocks)
  },

  // NB: BKD = best known digest
  [types.SET_RELAY_INFO] (state, { key, data }) {
    state.relay[key] = data
    lStorage.set('relay', state.relay)
  }
}

const actions = {
  // info: { source: String, date: Date }
  setLastComms ({ commit }, info) {
    commit(types.SET_LAST_COMMS, info)
  },

  // block: {
  //   height: Number,
  //   hash: String,
  //   verifiedAt: Date,
  //   isVerified: Boolean
  // }
  // Can pass one or all
  setCurrentBlock ({ commit }, block) {
    commit(types.SET_CURRENT_BLOCK, block)
  },

  async addPreviousBlock ({ commit }, previous) {
    return commit(types.ADD_PREVIOUS_BLOCK, previous)
  },

  // Called when there is a new current block
  // Relay-Info should trigger this in watch()
  async updateCurrentBlock ({ dispatch, state }, data) {
    await dispatch('addPreviousBlock', state.currentBlock)
    dispatch('setCurrentBlock', data)
  },

  // payload: { name: '', value: '' }
  setRelayInfo ({ commit }, payload) {
    commit(types.SET_RELAY_INFO, payload)
  },

  getExternalInfo ({ dispatch, state }) {
    console.log('Getting external info')
    axios.get(blockchainURL).then((res) => {
      console.log('EXTERNAL INFO:', res.data)
      const { height, hash } = res.data
      const currentHeight = state.currentBlock.height
      const currentHash = state.currentBlock.hash
      // NB: Do not change this weird spacing. Formats it pretty in console.
      console.log(`VERIFY:
Height:,
  Current:    ${currentHeight},
  New: ${height},
Digest:,
  Current:    ${currentHash},
  New: ${hash}
`)

      // If res.data.height > state.currentBlock.height, then verify and update
      if (height > currentHeight) {
        // Update current block
        dispatch('updateCurrentBlock', { height, hash, updatedAt: new Date() })
      }
      // Than verify height against relay
      dispatch('relay/verifyHeight', hash.toString(), { root: true })

      // Set last communication
      dispatch('setLastComms', { source: 'external', date: new Date() })
    }).catch((err) => {
      console.log('blockcypher error', err)
    })

  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
