
module.exports = {

  strip0xPrefix: function strip0xPrefix(hexString) {
    return hexString.substring(0, 2) === '0x' ? hexString.substring(2) : hexString;
  },

  concatenateHexStrings: function concatenateHexStrings(strs) {
    let current = '0x';
    for (let i = 0; i < strs.length; i += 1) {
      current = `${current}${this.strip0xPrefix(strs[i])}`;
    }
    return current;
  },
  concatenateHeadersHexes: function concatenateHeadersHexes(arr) {
    const hexes = arr.map(arr => arr.hex);
    return this.concatenateHexStrings(hexes);
  }
};
