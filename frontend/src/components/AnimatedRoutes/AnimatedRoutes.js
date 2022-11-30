import React from "react"
import { Route, Routes, useLocation } from "react-router-dom"
import Logs from "../Logs/Logs"
import NewProcess from "../NewProcess/NewProcess"
import ProcessList from "../ProcessList/ProcessList"
import Settings from "../Settings/Settings"
import {AnimatePresence} from "framer-motion"

function AnimatedRoutes() {
  const location = useLocation()
  return (
    <AnimatePresence>
      <Routes location={location} key={location.pathname}>
        <Route path="/" element={<ProcessList />} />
        <Route path="/new-process" element={<NewProcess />} />
        <Route path="/edit-process" element={<NewProcess />} />
        <Route path="/settings" element={<Settings />} />
        <Route path="/logs" element={<Logs />} />
      </Routes>
    </AnimatePresence>
  )
}

export default AnimatedRoutes
