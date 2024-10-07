export function MagicName(
  template: string,
  workDirname: string,
  title: string,
  index: number
): string {
  template = template.replace(/{{Title}}/g, title)
  template = template.replace(/{{Index}}/g, index.toString().padStart(3, '0'))
  template = template.replace(
    /{{RndInt}}/g,
    Math.floor(Math.random() * 1000)
      .toString()
      .padStart(3, '0')
  )
  template = template.replace(/{{RndChr}}/g, randomString(3))
  return template
}

function randomString(length: number): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

export function sanitizeFileName(input: string): string {
  // 创建一个正则表达式来匹配不允许的字符
  const re = /[<>:"/\\|?*\x00-\x1F]/g
  // 将不允许的字符替换为下划线
  let sanitized = input.replace(re, '_')

  // 去除首尾空白字符
  sanitized = sanitized.trim()
  // 去除开头和结尾的点号
  sanitized = sanitized.replace(/^\.*|\.*$/g, '')

  // 确保文件名长度不超过255个字符
  if (sanitized.length > 255) {
    sanitized = sanitized.substring(0, 255)
  }

  return sanitized
}
