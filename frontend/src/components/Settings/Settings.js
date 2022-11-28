import React, { Component, useContext } from "react"
import AppContext from "../../contexts/AppContext"

class Settings extends Component {
  constructor(props) {
    super(props)
    props.setConfig("Settings", true, "#00000000")
  }

  render() {
    return <div className="p-3">NOT IMPLEMENTED YET!</div>
  }
}

export default (props) => {
  return <Settings {...props} {...useContext(AppContext)} />
}
