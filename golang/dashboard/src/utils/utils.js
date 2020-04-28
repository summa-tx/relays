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
      return JSON.parse(i)
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

export function isString (str) {
  assert(str && typeof str === 'string', 'Must pass in string')
}

/**
 * Checks if string is hex
 * @param {String} str - string value to check
 * @returns {Boolean} true if string is hex, false if not
 */
export function isHex (str) {
  isString()

  if (/^[0-9a-fA-F]+$/.test(str)) {
    return true
  }

  return false
}

export function remove0x (str) {
  isString(str)

  if (str.slice(0, 2) === '0x') {
    return str.slice(2, str.length - 1)
  }
  return str
}

export function add0x (str) {
  isString(str)

  if (str.slice(0, 2) === '0x') {
    return str
  }
  return `0x${str}`
}

export function prependChars (chars, str) {
  isString(str)
  return `${chars}${str}`
}

export function reverseHex (str) {
  isHex(str)
  return str.match(/../g).reverse().join('')
}

export const convertEndian = {
  littleToBig (str) {
    isHex(str)
    const reversed = reverseHex(remove0x(str))
    return prependChars('00', reversed)
  },

  bigToLittle (str) {
    isHex(str)
    const reversed = reverseHex(str)
    return add0x(reversed)
  }
}

export function reverseEndianness (str) {
  if (str.slice(0,2) === "0x") {
    return str.slice(2).match(/../g).reverse().join("")
  }
  return str.match(/../g).reverse().join("")
}
