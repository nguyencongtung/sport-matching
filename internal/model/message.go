package model

     type Message struct {
         Room    string `json:"room"`
         Sender  string `json:"sender"`
         Content string `json:"content"`
     }