import { Add, PlayArrow, Stop, Delete, Edit } from "@material-ui/icons"
import { Component, useContext } from "react"
import { useNavigate } from "react-router-dom"
import "./ProcessList.css"
import "tw-elements"
import AppContext from "../../contexts/AppContext"
import { DragDropContext, Draggable, Droppable } from "react-beautiful-dnd"

class ProcessList extends Component {
  constructor(props) {
    super(props)
    props.setConfig("My Process Manager", false, "#00000000")
    this.onEvent = this.onEvent.bind(this)
    this.onDragEnd = this.onDragEnd.bind(this)
    this.state = {
      processes: [],
    }
  }

  onEvent(id, status) {
    console.log("process status >>>> " + id, status)
    this.setState((st) => {
      return {
        processes: st.processes.map((it) => {
          if (it.id === id) {
            return {
              ...it,
              run_status: status, // idle, running, error
            }
          }
          return it
        }),
      }
    })
  }

  componentDidMount() {
    if (!window.go) return
    ;(async () => {
      const list = await window.go.main.App.GetProcesses()
      console.log(list)

      this.setState({
        processes: list,
      })

      if (!window.runtime) return
      window.runtime.EventsOn("run-status", this.onEvent)
    })()
  }

  componentWillUnmount() {
    if (!window.runtime) return
    window.runtime.EventsOff("run-status")
  }

  componentDidUpdate() {
    this.props.updateTooltips()
  }

  onDragEnd(droppedItem) {
    if (!droppedItem.destination) return
    var updatedList = [...this.state.processes]
    const [reorderedItem] = updatedList.splice(droppedItem.source.index, 1)
    updatedList.splice(droppedItem.destination.index, 0, reorderedItem)

    this.setState({
      processes: updatedList
    })

    ;(async () => {
      const result = await window.go.main.App.ProcessesReorder(
        updatedList.map((it) => it.id)
      )
      console.log(result)
    })()
  }

  render() {
    const eachItem = (item) => {
      return (
        <div
          className="process-item cursor-pointer border-b-2 text-lg grid items-center justify-center pl-1"
          onClick={(e) => {
            this.props.navigate("/logs", {
              state: {
                id: item.id,
                name: item.name,
              },
            })
          }}
        >
          <div
            className="scale-50 px-1 pt-4"
            style={{
              fill:
                item.status === 0 || item.run_status !== "running"
                  ? "#c2c2c2"
                  : "#54d1b0",
            }}
          >
            <svg viewBox="0 0 350 239.2">
              <path
                d="M132,107.2c3.3,3.3,5.2,7.8,5.2,12.4s-1.9,9.1-5.2,12.4l-102,102c-4.4,4.4-10.9,6.2-16.9,4.6
                  c-6-1.6-10.8-6.3-12.4-12.3c-1.6-6,0.1-12.5,4.5-16.9l89.8-89.8L5.1,29.9C1.8,26.6,0,22.1,0,17.5c0-4.6,1.8-9.1,5.1-12.3
                  C8.4,1.9,12.8,0,17.5,0s9.1,1.9,12.4,5.2L132,107.2z M332.5,204.3H180.8c-6.2,0-12,3.3-15.2,8.8c-3.1,5.4-3.1,12.1,0,17.5
                  c3.1,5.4,8.9,8.8,15.2,8.8h151.7c6.3,0,12-3.3,15.2-8.8s3.1-12.1,0-17.5C344.5,207.7,338.8,204.3,332.5,204.3z"
              />
            </svg>
          </div>
          <div className="my-4 w-[90%]">
            <div className="mb-[0.125rem] truncate">{item.name}</div>
            <div className="flex justify-start space-x-2 items-center text-xs text-slate-600">
              {item.status === 0 ? (
                <div>DISABLED</div>
              ) : item.run_status === "running" ? (
                <div>RUNNING</div>
              ) : item.run_status === "error" ? (
                <div>ERROR</div>
              ) : (
                <div>IDLE</div>
              )}
            </div>
          </div>
          <div className="process-icons">
            <div
              data-bs-toggle="tooltip"
              data-bs-placement="bottom"
              title={item.status === 0 ? "RUN" : "STOP"}
              className="process-icon"
              onClick={(e) => {
                e.stopPropagation()
                if (!window.go) return
                ;(async () => {
                  try {
                    if (item.status === 0) {
                      await window.go.main.App.RunProcess(item.id)
                    } else {
                      await window.go.main.App.StopProcess(item.id)
                    }
                    this.setState((st) => {
                      return {
                        processes: st.processes.map((it) => {
                          if (it.id === item.id) {
                            return {
                              ...it,
                              status: 1 - it.status,
                            }
                          }
                          return it
                        }),
                      }
                    })
                  } catch (err) {
                    console.log(">>> ", err)
                  }
                })()
              }}
            >
              {item.status === 0 ? (
                <PlayArrow style={{ fontSize: "1.3rem" }} />
              ) : (
                <Stop style={{ fontSize: "1.3rem" }} />
              )}
            </div>
            <div
              data-bs-toggle="tooltip"
              data-bs-placement="bottom"
              title="DUPLICATE"
              className="process-icon"
              onClick={(e) => {
                e.stopPropagation()
                this.props.navigate("/new-process", {
                  state: {
                    ...item,
                    name: "Copy of " + item.name,
                  },
                })
              }}
            >
              <Add style={{ fontSize: "1.6rem" }} />
            </div>
            <div
              data-bs-toggle="tooltip"
              data-bs-placement="bottom"
              title="EDIT"
              className="process-icon"
              onClick={(e) => {
                e.stopPropagation()
                this.props.navigate("/edit-process", { state: item })
              }}
            >
              <Edit style={{ fontSize: "1.3rem" }} />
            </div>
            <div
              data-bs-toggle="tooltip"
              data-bs-placement="bottom"
              title="DELETE"
              offset={1000}
              className="process-icon"
              onClick={(e) => {
                e.stopPropagation()
                if (!window.go) return
                ;(async () => {
                  try {
                    await window.go.main.App.DeleteProcess(item.id)
                    this.setState((st) => {
                      return {
                        processes: st.processes.filter((it) => {
                          return it.id !== item.id
                        }),
                      }
                    })
                  } catch (e) {
                    console.log(e)
                  }
                })()
              }}
            >
              <Delete style={{ fontSize: "1.3rem" }} />
            </div>
          </div>
        </div>
      )
    }
    return (
      <DragDropContext onDragEnd={this.onDragEnd}>
        <Droppable droppableId="droppable">
          {(provided, snapshot) => (
            <div {...provided.droppableProps} ref={provided.innerRef}>
              {this.state.processes.map((item, index) => (
                <Draggable
                  key={item.id}
                  draggableId={"item-" + item.id}
                  index={index}
                >
                  {(provided, snapshot) => (
                    <div
                      ref={provided.innerRef}
                      {...provided.draggableProps}
                      {...provided.dragHandleProps}
                    >
                      {eachItem(item)}
                    </div>
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </div>
          )}
        </Droppable>
        <div
          data-bs-toggle="tooltip"
          title="NEW PROCESS"
          data-bs-placement="left"
          className="btn-floating absolute right-8 bottom-8"
          onClick={(e) => {
            e.stopPropagation()
            this.props.navigate("/new-process")
          }}
        >
          <Add style={{ fontSize: "3rem" }} />
        </div>
      </DragDropContext>
    )
  }
}

export default (props) => {
  return (
    <ProcessList
      {...props}
      navigate={useNavigate()}
      {...useContext(AppContext)}
    />
  )
}
