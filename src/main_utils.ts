import {App} from 'vue'
import MButton from './components/MButton.vue'
import MLoading from './components/MLoading.vue'

export default function (app: App) {
  app.component('m-button', MButton)
  app.component('m-loading', MLoading)

  app.config.globalProperties.$e = () => (r: any) => {  // Request error handler
    if (r.headers.get('Content-Type') && r.headers.get('Content-Type').indexOf('application/json') !== -1) {
      alert(r.body.description)
    } else {
      alert(r.body)
    }
  }
}
