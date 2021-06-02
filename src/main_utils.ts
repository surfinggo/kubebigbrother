import {App} from 'vue'
import MButton from './components/MButton.vue'
import MLoading from './components/MLoading.vue'

// import {DateTime} from 'luxon'

export default function (app: App) {
  // app.filter("lux", (t, format) => DateTime.fromRFC2822(t).toFormat(format || "yyyy-MM-dd HH:mm:ss"))

  app.component('m-button', MButton)
  app.component('m-loading', MLoading)

  app.config.globalProperties.$e = () => r => {  // Request error handler
    if (r.headers.get('Content-Type') && r.headers.get('Content-Type').indexOf('application/json') !== -1) {
      alert(r.body.description)
    } else {
      alert(r.body)
    }
  }
}
