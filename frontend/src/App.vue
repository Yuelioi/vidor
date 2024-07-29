<template>
    <main class="flex w-screen h-screen bg-base-200">
        <div class="left bg-slate-900 text-slate-200">
            <!-- 主界面左侧栏 -->
            <SideBar></SideBar>
        </div>
        <div class="right flex flex-col flex-1 h-full">
            <!-- 主界面顶部 -->
            <HeaderArea class="h-12 shadow-xl" />
            <!-- 主界面内容区域 -->
            <div class="flex-1 overflow-hidden">
                <router-view></router-view>
            </div>
        </div>
        <FooterArea />
    </main>
</template>
