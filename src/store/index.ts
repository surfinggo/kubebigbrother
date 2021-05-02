import {createStore} from 'vuex'

import app from './module_app'

export default createStore({
    modules: {
        app: app,
    }
})