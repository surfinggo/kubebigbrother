import {App} from 'vue'
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome'
import {library} from '@fortawesome/fontawesome-svg-core'
import {
  faFileAlt,
  faHandHolding,
  faSearch,
  faSpinner,
  faTachometerAlt,
  faTrash,
  faUnlink
} from '@fortawesome/free-solid-svg-icons'

library.add(faHandHolding, faFileAlt, faSearch, faSpinner, faTachometerAlt, faTrash, faUnlink)

export default function (app: App) {
  app.component('font-awesome-icon', FontAwesomeIcon)
}