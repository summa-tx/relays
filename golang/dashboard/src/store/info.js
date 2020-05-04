import axios from 'axios'
import * as types from '@/store/mutation-types'
import { lStorage, convertUnixTimestamp } from '@/utils/utils'

const state = {
  source: 'blockstream.info',
  // When was the last time a successful communication was made?
  // If last successful check was > 5 minutes ago, show flag
  lastComms: lStorage.get('lastCommsExternal') || undefined, // Date

  currentBlock: lStorage.get('currentBlock') || {
    height: 0,              // Number - Current block height, from external
    hash: '',               // String - Current block hash, from external
    time: undefined,        // Date - Current block timestamp, from external
    updatedAt: undefined,   // Date - When was this data updated
  },

  // Keep track of previous block information
  // If incoming block number increments, then move currentBlock info to here
  // and incoming block info goes to currentBlock
  previousBlocks: lStorage.get('previousBlocks') || []
}

const mutations = {
  [types.SET_LAST_COMMS] (state, date) {
    state.lastComms = date
    lStorage.set('lastCommsExternal', state.lastComms)
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
  }
}

const actions = {
  setLastComms ({ commit }, { date }) {
    commit(types.SET_LAST_COMMS, date)
  },

  // block: {
  //   height: Number,
  //   hash: String,
  //   updatedAt: Date,
  //   isVerified: Boolean
  // }
  // Can pass one or all
  setCurrentBlock ({ commit }, block) {
    commit(types.SET_CURRENT_BLOCK, block)
  },

  async addPreviousBlock ({ commit, state }, previous) {
    if (state.currentBlock.height > previous.height) {
      console.log('adding previous block')
      return commit(types.ADD_PREVIOUS_BLOCK, previous)
    }
    return
  },

  // Called when there is a new current block
  // Relay-Info should trigger this in watch()
  async updateCurrentBlock ({ dispatch, state }, data) {
    await dispatch('addPreviousBlock', state.currentBlock)
    dispatch('setCurrentBlock', data)
  },

  getExternalInfo ({ dispatch, state, rootState }) {
    console.log('Getting external info')
    axios.get(`${rootState.blockchainURL}/blocks`).then((res) => {
      console.log('EXTERNAL INFO:', res.data[0])
      const { height, id: hash, timestamp } = res.data[0]
      const currentHeight = state.currentBlock.height
      const currentHash = state.currentBlock.hash
      const time = convertUnixTimestamp(timestamp)

      console.log(`VERIFY\n\tHeight:\n\t\tCurrent: ${currentHeight},\n\t\tNew: ${height},
        \n\tDigest:\n\t\tCurrent: ${currentHash},\n\t\tNew: ${hash}`)

      dispatch('updateCurrentBlock', {
        height,
        hash,
        time,
        updatedAt: new Date()
      })

      // Set last communication
      dispatch('setLastComms', { date: new Date() })
    }).catch((err) => {
      console.log('blockstream error', err)
    })

  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
