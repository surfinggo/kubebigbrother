import {App} from 'vue'
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome'
import {library} from '@fortawesome/fontawesome-svg-core'
import {faHandHolding, faSpinner, faTachometerAlt, faTrash} from '@fortawesome/free-solid-svg-icons'

library.add(faHandHolding, faSpinner, faTachometerAlt, faTrash)

export default function (app: App) {
  app.component('font-awesome-icon', FontAwesomeIcon)
}