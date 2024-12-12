const { createApp } = Vue;

import Header from './components/Header.js';  // Import Header component
import DumpList from './components/DumpList.js';  // Import Counter component
import Counter from './components/Counter.js';  // Import Counter component

const App = {
    components: {
        Header,
        Counter,
        DumpList,
    },
    template: `
        <div>
            <h4>Database &amp; Volume Backups</h4>
            <Suspense>
                <DumpList/>
            </Suspense>
        </div>
    `
};

// Mount Vue app
createApp(App).mount('#app');
