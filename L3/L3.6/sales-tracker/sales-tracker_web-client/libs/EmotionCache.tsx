"use client"

import { CacheProvider } from "@emotion/react"
import createCache from "@emotion/cache"
import { useServerInsertedHTML } from "next/navigation"
import { useState } from "react"

export default function EmotionCacheProvider({ children }: { children: React.ReactNode }) {
  const [cache] = useState(() => {
    const cache = createCache({ key: "mui", prepend: true })
    cache.compat = true
    return cache
  })

  useServerInsertedHTML(() => {
    const styles = cache.inserted
    const names = Object.keys(styles)

    if (names.length === 0) {
      return null
    }

    let css = ""
    names.forEach((name) => {
      if (name !== "global") {
        css += styles[name]
      }
    })

    return (
      <style
        data-emotion={`${cache.key} ${names.join(" ")}`}
        dangerouslySetInnerHTML={{ __html: css }}
      />
    )
  })

  return <CacheProvider value={cache}>{children}</CacheProvider>
}