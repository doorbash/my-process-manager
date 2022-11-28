import "./App.css"

import Header from "./components/Header/Header"
import ProcessList from "./components/ProcessList/ProcessList"
import NewProcess from "./components/NewProcess/NewProcess"
import { Routes, Route, BrowserRouter as Router } from "react-router-dom"
import Settings from "./components/Settings/Settings"
import Logs from "./components/Logs/Logs"

import React, { Component } from "react"
import AppContext from "./contexts/AppContext"

export default class App extends Component {
  constructor(props) {
    super(props)
    this.setConfig = this.setConfig.bind(this)
    this.updateTooltips = this.updateTooltips.bind(this)
    this.state = {
      title: "",
      showBackButton: false,
      bgColor: "white",
    }
  }

  setConfig(title, showBackButton, bgColor) {
    this.setState({
      title: title,
      showBackButton: showBackButton,
      bgColor: bgColor,
    })
  }

  updateTooltips() {
    if (this.tooltipList) {
      this.tooltipList.forEach((tooltip) => {
        tooltip.dispose()
      })
    }

    this.tooltipList = [].slice
      .call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
      .map(function (tooltipTriggerEl) {
        return new window.Tooltip(tooltipTriggerEl)
      })

    this.tooltipList.forEach((tooltip) => {
      tooltip._config = {
        ...tooltip._config,
        fallbackPlacements: ["bottom"],
        offset: "0,15",
      }
      tooltip
        .getTipElement()
        .querySelector(".tooltip-inner")
        .classList.add("custom-tooltip")
    })
  }

  componentWillUnmount() {
    if (this.tooltipList) {
      this.tooltipList.forEach((tooltip) => {
        tooltip.dispose()
      })
    }
  }

  render() {
    return (
      <div className="app">
        <Router>
          <AppContext.Provider
            value={{
              ...this.state,
              setConfig: this.setConfig,
              updateTooltips: this.updateTooltips,
            }}
          >
            <Header />
            <div
              className="app-content"
              style={{ backgroundColor: this.state.bgColor }}
            >
              <Routes>
                <Route path="/" element={<ProcessList />} />
                <Route path="/new-process" element={<NewProcess />} />
                <Route path="/edit-process" element={<NewProcess />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="/logs" element={<Logs />} />
              </Routes>
            </div>
          </AppContext.Provider>
        </Router>
      </div>
    )
  }
}
