import DefaultTheme from 'vitepress/theme'
import Playground from '../components/Playground.vue'
import PageToc from '../components/PageToc.vue'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('Playground', Playground)
    app.component('PageToc', PageToc)
  }
}
