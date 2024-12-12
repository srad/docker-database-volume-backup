<template>
    <table class="table table-bordered table-hover">
        <thead>
        <tr>
            <th class="bg-light text-end" style="width: 10%">#</th>
            <th class="bg-light">Filename</th>
            <th class="bg-light" style="width: 10%">Size</th>
            <th class="bg-light" style="width: 20%">Created</th>
            <th class="bg-light text-center" style="width: 5%">Restore</th>
            <th class="bg-light text-center" style="width: 5%">Delete</th>
        </tr>
        </thead>
        <tbody>
        <tr class="align-middle" :key="file.filename" v-for="(file, i) in files">
            <td class="text-end">{{ i + 1 }}</td>
            <td>
                <a :href="'/backups/dumps/' + file.filename">{{ file.filename }}</a></td>
            <td class="text-end">{{ file.fileSize }}MB</td>
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

const formatDate = (goTimeString: string) => {
    const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    const date = new Date(goTimeString);

    // Extract parts
    const weekDay = date.getDay();
    const day = date.getDate(); // Day of the month (1-31)
    const month = date.getMonth() + 1; // Month (0-11, so add 1)
    const year = date.getFullYear(); // Year
    const hour = date.getHours(); // Hours (0-23)
    const minute = date.getMinutes(); // Minutes (0-59)

    return `${days[weekDay]}, ${day}.${month}.${year} ${hour}:${minute}`;
};

type Dump = {
    filename: string;
    fileSize: number;
    created: string;
}

const response = await client?.get<Dump[]>("/dumps");
const list = response?.data.map(x => {
    return {
        ...x,
        fileSize: (x.fileSize / 1024 / 1024).toFixed(2),
        created: formatDate(x.created)
    };
});
const files = ref(list);
</script>

<style scoped>
</style>
