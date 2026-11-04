import { defineStore } from 'pinia'
import { ref, watch, computed } from 'vue'

export type Theme = 'light' | 'dark'
export type UiStyle = 'cupertino' | 'blueprint'

const STORAGE_THEME_KEY = 'ems-theme'
const STORAGE_STYLE_KEY = 'eam-ui-style'
const DARK_CLASS = 'dark'

export const useThemeStore = defineStore('theme', () => {
  // 从 localStorage 读取主题与风格
  const storedTheme = localStorage.getItem(STORAGE_THEME_KEY) as Theme | null
  const storedStyle = localStorage.getItem(STORAGE_STYLE_KEY) as UiStyle | null

  const theme = ref<Theme>(storedTheme || 'light')
  const uiStyle = ref<UiStyle>(storedStyle || 'cupertino')

  // 应用主题 (light / dark)
  function applyTheme(themeValue: Theme) {
    const html = document.documentElement
    const body = document.body

    html.setAttribute('data-theme', themeValue)

    if (themeValue === 'dark') {
      html.classList.add(DARK_CLASS)
      body.classList.add(DARK_CLASS)
    } else {
      html.classList.remove(DARK_CLASS)
      body.classList.remove(DARK_CLASS)
    }
  }

  // 应用 UI 风格 (cupertino / blueprint)
  function applyUiStyle(styleValue: UiStyle) {
    const html = document.documentElement
    html.setAttribute('data-ui-style', styleValue)
  }

  // 初始化应用
  applyTheme(theme.value)
  applyUiStyle(uiStyle.value)

  // 监听主题变化
  watch(theme, (newTheme) => {
    applyTheme(newTheme)
    localStorage.setItem(STORAGE_THEME_KEY, newTheme)
  })

  // 监听风格变化
  watch(uiStyle, (newStyle) => {
    applyUiStyle(newStyle)
    localStorage.setItem(STORAGE_STYLE_KEY, newStyle)
  })

  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  function setTheme(newTheme: Theme) {
    theme.value = newTheme
  }

  function toggleUiStyle() {
    uiStyle.value = uiStyle.value === 'blueprint' ? 'cupertino' : 'blueprint'
  }

  function setUiStyle(newStyle: UiStyle) {
    uiStyle.value = newStyle
  }

  return {
    theme,
    uiStyle,
    isDark: computed(() => theme.value === 'dark'),
    isBlueprint: computed(() => uiStyle.value === 'blueprint'),
    toggleTheme,
    setTheme,
    toggleUiStyle,
    setUiStyle,
  }
})
