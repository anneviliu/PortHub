new Vue({
    el: '#app',
    delimiters: ["{[", "]}"],

    data() {
        return {
            res:[],
            portDetails: []
        }
    },
    methods: {
        async getResult() {
            var resp = (await axios.get('/api/getResult')).data;

            // Object copy
            this.portDetails = JSON.parse(JSON.stringify(resp.messages));
            
            for(let k in resp.messages) {
                var portCount = resp.messages[k].port.length
                var maxNum = 10
                if (portCount > maxNum) {
                    resp.messages[k].port = resp.messages[k].port.slice(0,maxNum)
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