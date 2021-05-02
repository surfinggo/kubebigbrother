export default {
  state: {
    siteName: '', // will be updated when app created
    siteDescription: '', // will be updated when app created
  },
  getters: {},
  mutations: {
    SET_SITE_NAME(state, payload) {
      state.siteName = payload
    },
    SET_SITE_DESCRIPTION(state, payload) {
      state.siteDescription = payload
    },
  },
  actions: {},
}
