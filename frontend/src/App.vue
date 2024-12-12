<template>
    <Suspense>
        <main class="container-fluid my-2">
            <h2 class="mb-3">Backups</h2>
            <BackupList/>
        </main>
    </Suspense>
</template>

<script setup lang="ts">
import BackupList from "./components/BackupList.vue";
</script>

<style scoped lang="scss">
</style>
