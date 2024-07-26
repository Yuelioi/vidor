export interface Tab {
    id: number
    name: string
    icon: string
    component: ReturnType<typeof defineComponent>
}
