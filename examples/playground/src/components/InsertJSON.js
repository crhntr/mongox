var InsertJSON = Vue.component('insert-json', {
  data () {
    return {
      raw: '{\n\t"text":"Hello, world!"\n}',
      doc: {},
      msg: '',
      textColor: '#000'
    }
  },
  props: ['col'],

  watch: {
    'raw': function () {
      try {
        if (this.raw) {
          this.doc = JSON.parse(this.raw)
        }

        this.msg = ''
        this.textColor = '#000'
      } catch (e) {
        this.textColor = '#F00'
        this.msg = 'json error'
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

  template: `<div class="insert-json">
  <textarea v-model="raw"
    :style="{color: textColor}"
    placeholder="insert document"
    onkeydown="if(event.keyCode===9){var v=this.value,s=this.selectionStart,e=this.selectionEnd;this.value=v.substring(0, s)+'\t'+v.substring(e);this.selectionStart=this.selectionEnd=s+1;return false;}"></textarea>
  <p v-if="msg">{{msg}}</p>
  <button @click="insert()">Insert into {{col}}</button>
</div>`
})
