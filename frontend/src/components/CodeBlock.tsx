import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import graphql from 'highlight.js/lib/languages/graphql'
import json from 'highlight.js/lib/languages/json'
import plaintext from 'highlight.js/lib/languages/plaintext'
import 'highlight.js/styles/atom-one-dark.css'
import { useEffect } from 'react'

interface Props {
  code: string
  language: string
  maxHeight?: number
}

export const CodeBlock = ({ code, language, maxHeight }: Props) => {
  useEffect(() => {
    hljs.configure({ ignoreUnescapedHTML: true })
    hljs.registerLanguage('graphql', graphql)
    hljs.registerLanguage('json', json)
    hljs.registerLanguage('go', go)
    hljs.registerLanguage('plaintext', plaintext)
    hljs.highlightAll()
  })

  return (
    <pre>
      <code className={`max-h-[${maxHeight}px] language-${language} text-xs`}>{code}</code>
    </pre>
  )
}
