import {createApp} from 'vue'

import router from './main_router'
import store from './store'
import axios from './main_axios'
import fontawesome from './main_fontawesome'
import utils from './main_utils'

import App from './App.vue'

import 'virtual:windi.css'

import './index.css'


const app = createApp(App)

app.use(router)
app.use(store)
app.use(axios)
app.use(fontawesome)
app.use(utils)
app.mount('#app')
