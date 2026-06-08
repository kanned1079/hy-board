export const renderMarkdown = (text: string | undefined | null): string => {
  if (!text) return ''
  
  // 1. HTML Escape to prevent XSS
  let html = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // 2. Headers
  html = html.replace(/^### (.*$)/gim, '<h4 class="text-sm font-extrabold mt-2 mb-1">$1</h4>')
  html = html.replace(/^## (.*$)/gim, '<h3 class="text-md font-extrabold mt-3 mb-1.5">$1</h3>')
  html = html.replace(/^# (.*$)/gim, '<h2 class="text-lg font-extrabold mt-4 mb-2">$1</h2>')

  // 3. Bold (**text**)
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')

  // 4. Italic (*text*)
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')

  // 5. Code blocks (```code```)
  html = html.replace(/```([\s\S]*?)```/g, '<pre class="bg-zinc-100 dark:bg-zinc-900 p-2 rounded text-[10px] font-mono my-2 overflow-x-auto select-all">$1</pre>')

  // 6. Inline code (`code`)
  html = html.replace(/`(.*?)`/g, '<code class="bg-zinc-100 dark:bg-zinc-900 px-1 py-0.5 rounded text-[10px] font-mono">$1</code>')

  // 7. Unordered Lists (- item)
  html = html.replace(/^\s*[-*]\s+(.*$)/gim, '<li class="ml-4 list-disc">$1</li>')

  // 8. Hyperlinks ([title](url))
  html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-primary-500 dark:text-primary-400 hover:underline font-medium">$1</a>')

  // 9. Process line breaks (only outside <pre> tags to preserve formatting)
  const parts = html.split(/(<pre[\s\S]*?<\/pre>)/g)
  for (let i = 0; i < parts.length; i++) {
    if (!parts[i].startsWith('<pre')) {
      parts[i] = parts[i].replace(/\n/g, '<br>')
    }
  }
  return parts.join('')
}
