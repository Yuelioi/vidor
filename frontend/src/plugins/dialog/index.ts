import type { App, Plugin } from 'vue'
import VDialog from './dialog.vue'

export { VDialog }

export const VDialogPlugin: Plugin = {
    install(app: App) {
        app.component('VDialog', VDialog)
    }
}
