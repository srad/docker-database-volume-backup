<template>
    <table class="table table-bordered table-hover">
        <thead>
        <tr>
            <th class="bg-light text-end" style="width: 10%">#</th>
            <th class="bg-light">Filename</th>
            <th class="bg-light" style="width: 20%">Created</th>
            <th class="bg-light text-center" style="width: 5%">Restore</th>
            <th class="bg-light text-center" style="width: 5%">Delete</th>
        </tr>
        </thead>
        <tbody>
        <tr class="align-middle" :key="file.backupId" v-for="(file, i) in files">
            <td class="text-end">{{ i + 1 }}</td>
            <td>
                <div class="mb-1">
                    <a :href="'/backupfiles/dumps/' + file.databaseDumpFilename"> {{ file.databaseDumpFilename }} </a>
                    <span class="form-text">({{ file.dbSize }} MB)</span>
                </div>

                <div>
                    <a :href="'/backupfiles/files/' + file.volumeDumpFilename">{{ file.volumeDumpFilename }} </a>
                    <span class="form-text">({{ file.volumeSize }} MB)</span>
                </div>
            </td>
            <td class="text-end">{{ file.created }}</td>
            <td class="text-center p-1">
                <button type="button" class="btn btn-sm btn-warning">Restore</button>
            </td>
            <td class="text-center p-1">
                <button type="button" class="btn btn-sm btn-danger">Delete</button>
            </td>
        </tr>
        </tbody>
    </table>
</template>

<script setup lang="ts">
import {ref} from 'vue';
import client from "../client.ts";
import {formatDate} from '../lib/datetime';

type Backup = {
    backupId: number;
    databaseDumpFilePath: string;
    databaseDumpFilename: string;
    databaseDumpFileSize: number;
    volumeDumpFilePath: string;
    volumeDumpFilename: string;
    volumeDumpFileSize: number;
    created: string;
}

const response = await client?.get<Backup[]>("/backups");
const list = response?.data.map(x => {
    return {
        ...x,
        dbSize: (x.databaseDumpFileSize / 1024 / 1024).toFixed(2),
        volumeSize: (x.volumeDumpFileSize / 1024 / 1024).toFixed(2),
        created: formatDate(x.created)
    };
});
const files = ref(list);
</script>

<style scoped>
</style>
