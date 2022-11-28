package main

type Log struct {
	Time    int64  `json:"time"`
	Type    string `json:"type"` // out, err
	Message string `json:"message"`
}
