module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  pluginOptions: {
    express: {
      shouldServeApp: true,
      serverDir: './server'
    }
  },
  devServer: {
    proxy: 'http://localhost:1317'
  }
}
