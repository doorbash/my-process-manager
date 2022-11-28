import React, { Component, useContext } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import AppContext from "../../contexts/AppContext"
import "./Logs.css"

const MAX_LINES = 50

class Logs extends Component {
  constructor(props) {
    super(props)
    props.setConfig(props.location.state.name, true, "black")
    this.onEvent = this.onEvent.bind(this)
    this.state = {
      logs: [],
    }
  }

  onEvent(id, time, type, message) {
    console.log("run log >>>>>>", id, time, type, message)
    if (id !== this.props.location.state.id) return
    this.setState((st) => {
      return {
        logs: [
          ...st.logs,
          {
            time: time,
            type: type,
            message: message,
          },
        ].slice(-MAX_LINES),
      }
    })
  }

  componentDidMount() {
    // this.timer = setInterval(() => {
    //   this.onEvent(0, null, "out", "New Log! " + Date.now())
    // }, 300)
    if (!window.go) return
    ;(async () => {
      const logs = await window.go.main.App.GetLogs(
        this.props.location.state.id
      )
      // console.log(logs)
      this.setState((st) => {
        return {
          logs: logs.slice(-MAX_LINES),
        }
      })
    })()
    if (!window.runtime) return
    window.runtime.EventsOn("run-log", this.onEvent)
  }

  componentWillUnmount() {
    // clearInterval(this.timer)
    if (!window.runtime) return
    window.runtime.EventsOff("run-log")
  }

  componentDidUpdate() {
    let div = document.querySelector(".app-content")
    // let isScrolledToBottom =
    //   div.scrollHeight - div.clientHeight <=
    //   div.scrollTop + div.offsetHeight * 0.25
    // if (isScrolledToBottom) {
      div.scrollTop = div.scrollHeight - div.clientHeight
    // }
  }

  render() {
    return (
      <div className="logs">
        {this.state.logs.map((item) => {
          return (
            <div className="log text-xs p-1">
              <p>{item.message}</p>
            </div>
          )
        })}
      </div>
    )
  }
}

export default (props) => {
  return (
    <Logs
      {...props}
      navigate={useNavigate()}
      location={useLocation()}
      {...useContext(AppContext)}
    />
  )
}
