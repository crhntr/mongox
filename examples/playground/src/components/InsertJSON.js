

var InsertJSON = Vue.component('insert-json', {
  data () {
    return {
      raw: '{\n\t"text":"Hello, world!"\n}',
      doc: {},
      msg: ''
    }
  },
  props: ['col'],

  watch: {
    'raw': function () {
      try {
        this.doc = JSON.parse(this.raw)
        this.msg = ''
      } catch (e) {
        this.msg = e
      }
    }
  },

  methods: {
    insert: function () {
      this.doc = JSON.parse(this.raw)
      this.$store.dispatch('db/'+this.col+'/insert', this.doc).then(this.resolve, this.reject)
    },
    resolve: function (res) {
      console.log(res)
      this.doc = {}
      this.raw = ''
      this.msg = "inserted"
    },
    reject: function (err) {
      this.doc = {}
      this.raw = ''
      this.msg = "Error " + JSON.stringify(err)
    }
  },

  template: `<div>
  <textarea v-model="raw"
    onkeydown="if(event.keyCode===9){var v=this.value,s=this.selectionStart,e=this.selectionEnd;this.value=v.substring(0, s)+'\t'+v.substring(e);this.selectionStart=this.selectionEnd=s+1;return false;}"></textarea>
  <p v-if="msg">{{msg}}</p>
  <button @click="insert()">Insert into {{col}}</button>
</div>`
})
