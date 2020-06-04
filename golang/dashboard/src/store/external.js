import axios from 'axios'
import * as types from '@/store/mutation-types'
import { lStorage, convertUnixTimestamp } from '@/utils/utils'

const state = {
  source: 'blockstream.info',

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
  [types.SET_LAST_COMMS_EXTERNAL] (state, date) {
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
  addPreviousBlock ({ commit, state }, newBlock) {
    if (newBlock.height > state.currentBlock.height) {
      commit(types.ADD_PREVIOUS_BLOCK, state.currentBlock)
    }
  },

  async updateCurrentBlock ({ dispatch, commit }, newBlock) {
    await dispatch('addPreviousBlock', newBlock)
    commit(types.SET_CURRENT_BLOCK, newBlock)

  },

  getExternalInfo ({ dispatch, commit, rootState }) {
    console.log('Getting external info')
    axios.get(`${rootState.blockchainURL}/blocks`).then((res) => {
      console.log('EXTERNAL INFO:', res.data[0])
      const { height, id: hash, timestamp } = res.data[0]
      const time = convertUnixTimestamp(timestamp)

      dispatch('updateCurrentBlock', {
        height,
        hash,
        time,
        updatedAt: new Date()
      })

      commit(types.SET_LAST_COMMS_EXTERNAL, new Date())
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
