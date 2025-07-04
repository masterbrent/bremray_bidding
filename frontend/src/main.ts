import { mount } from 'svelte'
import './styles/global.css'
import './app.css'
import App from './App.svelte'
// import App from './AppMinimal.svelte'

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
