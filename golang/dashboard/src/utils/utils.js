/**
 * How many minutes has passed from a past point until now?
 * Expected usage: current block verified at, last comms relay, last comms
 * external
 * @param {Date}      date Starting date
 * @returns {Number}  Minutes that have passed from starting date until now
 *
 */
export function getMinsAgo (date) {
  if (!date) {
    return undefined
  }
  const from = timeInSecs(date)
  const now = timeInSecs(new Date())
  return Math.round((now - from) / 60)
}

/**
 * Gets time in seconds from a date object
 * @param {Date}      date if no date is passed in, current date is used
 * @returns {Number}  Returns time in seconds
 */
export function timeInSecs (date) {
  const d = new Date(date) || new Date()
  return d.getTime() / 1000
}

/**
 * Convenience class to verify localStorage exists
 * TODO: Consider adding polyfill...
 *
 */
class LStorage {
  _verify () {
    if (window && window.localStorage) {
      return true
    }
    return false
  }

  /**
   * Sets item to localStorage
   * @param {String} item - name of item to set
   * @param {Any} value - value of item to set
   */
  set (item, value) {
    if (this._verify()) {
      window.localStorage.setItem(item, JSON.stringify(value))
    } else {
      console.error('Error saving value to localStorage')
    }
  }

  /**
   * Gets items from localStorage
   * @param {String} item - name of item to retrieve
   * @returns {String}
   */
  get (item) {
    if (this._verify()) {
      const i = window.localStorage.getItem(item)

      let value
      try {
        value = JSON.parse(i)
      } catch (e) {
        console.log('storage error', e)
      }
      return value
    } else {
      console.error('Error getting value from localStorage')
    }
  }

  /**
   * Removes an item from localStorage
   * @param {String} item - name of item to remove
   */
  remove (item) {
    if (this._verify()) {
      window.localStorage.removeItem(item)
    } else {
      console.error('Error removing value from localStorage')
    }
  }
}

export const lStorage = new LStorage()

const assert = require('bsert')

/**
 * Checks if value is of type string
 * @param {String} str - string value to check
 * @returns {Boolean} true if value is string, false if not
 */
export function isString (str) {
  const isStr = typeof str === 'string'
  assert(isStr, `Must pass in string, received ${typeof str}`)
}

/**
 * Checks if string is hex
 * @param {String} str - string value to check
 * @returns {Boolean} true if string is hex, false if not
 */
export function isHex (str) {
  isString(str)

  let hexStr = remove0x(str)

  assert(hexStr && /^[0-9a-fA-F]+$/.test(hexStr), 'Must pass in hex string')
}

/**
 * If a hex string is '0x' prepended, it removes it
 * @param {String} str - hex string
 * @returns {String} hex string without '0x'
 */
export function remove0x (str) {
  isString(str)

  if (str.slice(0, 2) === '0x') {
    return str.slice(2, str.length)
  }
  return str
}

/**
 * If a hex string is not already '0x' prepended, it adds it
 * @param {String} str - hex string
 * @returns {String} hex string beginning with '0x'
 */
export function add0x (str) {
  isString(str)

  if (str.slice(0, 2) === '0x') {
    return str
  }
  return `0x${str}`
}

/**
 * Reverses Endianness of a hex bytes string
 * @param {String} str - hex string
 * @returns {String} hex string with reverse endianness
 */
export function reverseEndianness (str) {
  var formatStr = remove0x(str)
  return formatStr.match(/../g).reverse().join('')
}

/**
 * Converts a Unix timestamp
 * BlockStream returns Unix timestamps that must be converted
 * @param {Number} time - hex string
 * @returns {Date} time as a JavaScript Date object
 */
export function convertUnixTimestamp (time) {
  return new Date(time * 1000)
}
