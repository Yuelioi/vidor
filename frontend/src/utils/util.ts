export function MagicName(
    template: string,
    workDirname: string,
    title: string,
    index: number
): string {
    template = template.replace(/{{WorkDir}}/g, workDirname)
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
