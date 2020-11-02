import App from "./js/components/App"
import ReactDOM from "react-dom";
import React from "react";
import "./css/main.css"

const wrapper = document.getElementById("container");
wrapper ? ReactDOM.render(<App />, wrapper) : false;