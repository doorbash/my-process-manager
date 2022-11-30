import React, { Component, useContext } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import "tw-elements"
import AppContext from "../../contexts/AppContext"
import "./NewProcess.css"
import { motion } from "framer-motion"

class NewProcess extends Component {
  constructor(props) {
    super(props)
    if (props.location.pathname === "/edit-process") {
      props.setConfig("Edit Process", true, "white")
    } else {
      props.setConfig("New Process", true, "white")
    }
    this.handleSubmit = this.handleSubmit.bind(this)

    if (props.location.state) {
      this.state = {
        ...props.location.state,
        error: undefined,
      }
    } else {
      this.state = {
        name: "My Process",
        command: "ping 8.8.8.8 -t",
        error: undefined,
      }
    }
  }

  handleSubmit(e) {
    e.preventDefault()
    console.log("form submitted!")
    this.setState({
      error: undefined,
    })
    if (!window.go) return
    ;(async () => {
      try {
        if (this.props.location.pathname === "/edit-process") {
          await window.go.main.App.UpdateProcess({
            ...this.state,
            create_time: Date.now(),
          })
          this.props.navigate(-1)
        } else {
          await window.go.main.App.InsertProcess({
            ...this.state,
            create_time: Date.now(),
            status: 1,
          })
          this.props.navigate(-1)
        }
      } catch (e) {
        this.setState({
          error: e,
        })
      }
    })()
  }

  render() {
    return (
      <motion.div
        initial={{
          position: "absolute",
          left: 0,
          right: 0,
          opacity: 0,
          transition: { duration: 0 },
        }}
        animate={{
          position: "static",
          opacity: 1,
          transition: { duration: 0.3 },
        }}
        exit={{
          position: "absolute",
          left: 0,
          right: 0,
          opacity: 0,
          transition: { duration: 0 },
        }}
      >
        <form
          className="px-8 pt-4 mb-4 grid grid-cols-4 items-center gap-3 place-items-start"
          onSubmit={this.handleSubmit}
          autoComplete="off"
          spellCheck="false"
        >
          <label className="new-process-label" htmlFor="process-name">
            Name:
          </label>
          <input
            className="new-process-input"
            id="process-name"
            required={true}
            defaultValue={this.state.name}
            autoComplete="off"
            aria-autocomplete="none"
            spellCheck="false"
            onChange={(e) => {
              this.setState({ name: e.target.value })
            }}
          ></input>
          <textarea
            className="new-process-command"
            id="command"
            rows={8}
            required={true}
            defaultValue={this.state.command}
            autoComplete="off"
            aria-autocomplete="none"
            spellCheck="false"
            onChange={(e) => {
              this.setState({ command: e.target.value })
            }}
          ></textarea>
          <div className="new-process-submit">
            <input className="btn-primary" type="submit" value="save" />
          </div>
        </form>
        {this.state.error && (
          <div
            className="error-alert alert bg-red-100 rounded-lg py-5 px-6 mb-1 text-base text-red-700 flex items-center w-full fixed right-0 left-0 bottom-0"
            role="alert"
          >
            <strong className="mr-1">Error</strong> while inserting to database:{" "}
            {this.state.error}
          </div>
        )}
      </motion.div>
    )
  }
}

export default (props) => {
  return (
    <NewProcess
      {...props}
      navigate={useNavigate()}
      location={useLocation()}
      {...useContext(AppContext)}
    />
  )
}
