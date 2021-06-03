import {App} from 'vue'
// @ts-ignore
import Prism from 'vue-prism-component'

import 'prismjs/components/prism-json'
import 'prismjs/components/prism-yaml'
import 'prismjs/themes/prism-tomorrow.css'

export default function (app: App) {
  app.component('prism', Prism)
}