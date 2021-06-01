import {App} from 'vue'
import VueAxios from 'vue-axios'
import axios from 'axios'

export default function (app: App) {
  app.use(VueAxios, axios)
}
