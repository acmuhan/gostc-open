<script setup>
import {computed, h, nextTick, ref, watch} from "vue";
import MenuIconSvg from "./MenuIconSvg.vue";
import router, {allRouters} from "../../router/index.js";
import {localStore} from "../../store/local.js";
import {appStore} from "../../store/app.js";
import {routerToMenu} from "../../utils/routerToMenu.js";
import {normalRouters} from "../../router/routers/normal.js";

const menuTreeDataComputed = computed(() => {
  var funcMap = new Map()
  funcMap.set("funcWeb",appStore().siteConfig.funcWeb)
  funcMap.set("funcForward",appStore().siteConfig.funcForward)
  funcMap.set("funcTunnel",appStore().siteConfig.funcTunnel)
  funcMap.set("funcP2P",appStore().siteConfig.funcP2P)
  funcMap.set("funcProxy",appStore().siteConfig.funcProxy)
  funcMap.set("funcTun",appStore().siteConfig.funcTun)
  funcMap.set("funcNode",appStore().siteConfig.funcNode)
  funcMap.set("funcCfg",appStore().siteConfig.funcCfg)
  if (appStore().userInfo.admin === 1) {
    return routerToMenu(allRouters,funcMap)
  }
  return routerToMenu(normalRouters,funcMap)
})

const renderMenuIcon = (option) => {
  return h(MenuIconSvg, {
    svg: option.iconSvg
  }, null);
}

const menuSelectChange = (key, item) => {
  if (item?.link) {
    window.open(item.link, '_blank')
    localStore().isCollapsed = true
    return
  }
  if (router.currentRoute.value.name === key) {
    localStore().isCollapsed = true
    return
  }
  localStore().menuKey = key
  router.push({
    name: key
  })
  localStore().isCollapsed = true
}

const onMaskClick = () => {
  localStore().isCollapsed = true
}
const menu = ref()
watch(router.currentRoute, () => {
  nextTick(() => {
    menu.value?.showOption()
  })
})


</script>

<template>
  <div class="side-menu-container">
    <n-drawer :show="!localStore().isCollapsed"
              width="80%"
              placement="left"
              :trap-focus="false"
              :block-scroll="true"
              :on-mask-click="onMaskClick"
    >
      <n-drawer-content :title="appStore().siteConfig.title">
        <n-scrollbar style="height: calc(100vh - 60px)">
          <n-menu
              accordion
              ref="menu"
              :options="menuTreeDataComputed"
              :render-icon="renderMenuIcon"
              :on-update:value="menuSelectChange"
              :value="localStore().menuKey"
          />
        </n-scrollbar>
      </n-drawer-content>
    </n-drawer>
  </div>

</template>

<style scoped lang="scss">
.side-menu-container {
  position: fixed;
  inset: 0;
  z-index: 3000;
  width: 0;
  height: 0;
}

:deep(.n-drawer-mask),
:deep(.n-drawer) {
  position: fixed !important;
}

:deep(.n-drawer .n-drawer-content) {
  max-width: 320px;
}

:deep(.n-menu-item-content--selected)::before {

}

//:deep(.n-menu-item-content--child-active):not(.n-menu-item-content--collapsed)::before {
:deep(.n-menu-item-content--child-active.n-menu-item-content--collapsed)::before {

}

:deep(.n-drawer .n-drawer-content .n-drawer-body-content-wrapper) {
  box-sizing: border-box;
  padding: 0 !important;
}

:deep(.n-drawer-header) {
  height: 60px !important;
  box-sizing: border-box !important;
}
</style>