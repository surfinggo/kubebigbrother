import {createApp} from 'vue'
import VueAxios from 'vue-axios'
import axios from 'axios'
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome'
import {library} from '@fortawesome/fontawesome-svg-core'
import {faHandHolding, faSpinner, faTachometerAlt, faTrash} from '@fortawesome/free-solid-svg-icons'

import router from './router'
import store from './store'
import utils from './utils'

import App from './App.vue'

import './index.css'

library.add(faHandHolding, faSpinner, faTachometerAlt, faTrash)

const app = createApp(App)

app.use(router)
app.use(store)
app.use(utils)
app.use(VueAxios, axios)

app.component('fa', FontAwesomeIcon)

app.mount('#app')
