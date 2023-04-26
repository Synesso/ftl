export default {
  data() {
    return {
      activeVerb: null,
      message: "Starting warp drive...",
      moduleName: "com.squareup.ftldemo",
      verbs: [
        { key: 1, name: "NotifyVerb" },
        { key: 2, name: "PaymentVerb" },
        { key: 3, name: "PizzaVerb" },
      ],
    }
  },
  methods: {
    zoomIn(event, verb) {
      this.activeVerb = verb.key;
      console.log("selected: ", verb, event.target)
    }
  }
}
