new Vue({
    el: '#app',
    delimiters: ["{[", "]}"],

    data() {
        return {
            ips:[],
        }
    },
    methods: {
        async getResult() {
            var t = (await axios.get('/api/getResult')).data.messages;
            this.ips = t
            console.log(t)
        },
    },
    mounted () {
        this.getResult()
    }


});