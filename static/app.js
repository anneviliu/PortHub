new Vue({
    el: '#app',
    delimiters: ["{[", "]}"],

    data() {
        return {
            res:[],
        }
    },
    methods: {
        async getResult() {
            var resp = (await axios.get('/api/getResult')).data;

            for(let k in resp.messages) {
                var portCount = resp.messages[k].port.length
                if (portCount > 6) {
                    resp.messages[k].port = resp.messages[k].port.slice(0,7)
                }
            }
            this.res = resp
        },
    },
    mounted () {
        this.getResult();
        window.setInterval(() => {
            setTimeout(this.getResult(),0)
        }, 1000)
    }
});