new Vue({
    el: '#app',
    delimiters: ["{[", "]}"],

    data() {
        return {
            ips:[],
        }
    },
    methods: {
        getResult() {
            axios.get('/api/getResult')
                .then(function (response) {
                    // console.log(response.data.messages);
                    this.ips = response.data;
                    // if (response.data.msg == "登录成功"){
                    //     alert(response.data.msg);
                    //     location.href= "/"
                    // }
                    // else {
                    //     alert(response.data.msg);
                    // }
                })
                .catch(function (error) {
                    console.log(error);
                });
        }
    },
    mounted () {
        this.getResult()
        console.log(this.ips)
    }


});