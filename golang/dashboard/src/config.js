/** Config Variables */

// NODE
const NODE_ENV = process.env.NODE_ENV

// VUE APP. Must be prepended with VUE_APP in .env file for app to recognize
const NET_TYPE = process.env.VUE_APP_NET_TYPE
const CHAIN_NET = process.env.VUE_APP_CHAIN_NET
const RELAY_ADDRESS = process.env.VUE_APP_RELAY_ADDRESS
const DEBUG_BUTTONS = process.env.VUE_APP_DEBUG_BUTTONS

module.exports = {
  netType: NET_TYPE || 'local',
  chainNet: CHAIN_NET || 'cosmos',
  relayAddress: RELAY_ADDRESS || '1317',
  isProd: NODE_ENV === 'production',
  debugButtons: DEBUG_BUTTONS === 'true'
}
