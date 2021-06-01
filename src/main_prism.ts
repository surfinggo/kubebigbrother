import {App} from 'vue'
// @ts-ignore
import Prism from 'vue-prism-component'

import 'prismjs/themes/prism-tomorrow.css'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-yaml'

export default function (app: App) {
  app.component('prism', Prism)
}