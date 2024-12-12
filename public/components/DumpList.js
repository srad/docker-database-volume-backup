const template = `
    <table className="table table-bordered table-hover">
        <thead>
        <tr>
            <th className="bg-light">#</td>
            <th className="bg-light">Filename</td>
            <th className="bg-light" style="width: 10%">Size</td>
            <th className="bg-light" style="width: 10%">Created</td>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(file, i) in files">
            <td class="text-end">{{i + 1}}</td>
            <td>
                <a href="/static/{{file.filename}}">{{file.filename}}</a></td>
            <td class="text-end">{{(file.fileSize / 1024 / 1024).toFixed(2)}}MB</td>
            <td class="text-end">{{formatDate(file.created)}}</td>
        </tr>
        </tbody>
    </table>`;


const formatDate = (goTimeString) => {
    const date = new Date(goTimeString);
    const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

// Extract parts
    const weekDay = date.getDay();
    const day = date.getDate(); // Day of the month (1-31)
    const month = date.getMonth() + 1; // Month (0-11, so add 1)
    const year = date.getFullYear(); // Year
    const hour = date.getHours(); // Hours (0-23)
    const minute = date.getMinutes(); // Minutes (0-59)

    return `${days[weekDay]}, ${day}.${month}.${year} ${hour}:${minute}`;
}

export default {
    name: 'DumpList',
    async setup() {
        const {ref} = Vue;
        const dumps = await fetch("/api/dumps");
        const data = await dumps.json();
        const files = ref(data);

        return {
            files,
            formatDate,
        };
    },
    template: template,
};
