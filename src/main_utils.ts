import {App} from 'vue'
// import Block from './components/Block'
import MButton from './components/MButton.vue'
import MLoading from './components/MLoading.vue'
// import DefaultHeader from './components/DefaultHeader'
// import DefaultPagination from './components/DefaultPagination'
//
// import DefaultAvatar from './assets/images/default_avatar.png'
//
// import {DateTime} from 'luxon'
//
export default function (app: App) {
//   Vue.filter("lux", (t, format) => DateTime.fromRFC2822(t).toFormat(format || "yyyy-MM-dd HH:mm:ss"))
//
//   Vue.component('block', Block)
//   Vue.component('default-header', DefaultHeader)
//   Vue.component('default-pagination', DefaultPagination)
    app.component('m-button', MButton)
    app.component('m-loading', MLoading)
//
//   Vue.prototype.$defaultAvatar = DefaultAvatar
//
    app.config.globalProperties.$e = () => r => {  // Request error handler
        if (r.headers.get('Content-Type') && r.headers.get('Content-Type').indexOf('application/json') !== -1) {
            alert(r.body.description)
        } else {
            alert(r.body)
        }
    }
//
//   Vue.prototype.$refresh = function () {
//     this.$store.commit("SET_SITE_NAME", "GCS Browser")
//     this.$store.commit("SET_SITE_DESCRIPTION", "")
//   }
}
