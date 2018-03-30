var DataBrowser = Vue.component('data-browser', {
  props: ['col'],
  data () {
    return {
      query: '',
      projection: '',
      doc: null
    }
  },

  methods: {
    find: function () {
      this.$store.dispatch('db/'+this.col+'/find', {
        query: this.query,
        projection: this.projection || undefined
      }).then(this.resolve, this.reject)
    },
    resolve: function (res) {
      this.doc = res.body
    },
    reject: function () {
      this.doc = null
    }
  },

  template: `<div class="data-browser">
    <textarea placeholder="query"
      v-model="query"
      onkeydown="if(event.keyCode===9){var v=this.value,s=this.selectionStart,e=this.selectionEnd;this.value=v.substring(0, s)+'\t'+v.substring(e);this.selectionStart=this.selectionEnd=s+1;return false;}"/>
    </textarea>
    <textarea placeholder="projection"
      v-model="projection"
      onkeydown="if(event.keyCode===9){var v=this.value,s=this.selectionStart,e=this.selectionEnd;this.value=v.substring(0, s)+'\t'+v.substring(e);this.selectionStart=this.selectionEnd=s+1;return false;}"/>
    </textarea>

    <button @click="find">find</button>

    <div class="json">
      <div v-for="(val, name) in doc" :json-key="name" :key="key._id">
        {{val}}
      </div>
    </div>
  </div>`
})
