/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'


// import { VSnackbar } from "vuetify/lib/components/VSnackbar/VSnackbar.mjs";
// import { VBtn } from "vuetify/lib/components/VBtn/VBtn.mjs";
// import { VIcon } from "vuetify/lib/components/VIcon/VIcon.mjs";
// import { VDataTable } from "vuetify/lib/components/VDataTable/VDataTable.mjs";
// import { VStepper } from "vuetify/lib/components/VStepper/VStepper.mjs";
// import { VDatePicker } from "vuetify/lib/components/VDatePicker/VDatePicker.mjs";

export default createVuetify({
  theme: {
    themes: {
      light: {
        colors: {
          primary: '#1867C0',
          secondary: '#5CBBF6',
        },
      },
      dark: {
        colors: {
          primary: '#1867C0',
          secondary: '#5CBBF6',
        },
      },
    },
  },
  components: {
    // VDataTable,
    // VStepper, 
    // VDatePicker,
    // VSnackbar,
  }
})
