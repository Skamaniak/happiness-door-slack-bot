import React, {Component} from "react";

import {MessageType} from "../api/Protocol"

import Socket from '../api/Socket';
import LoadingIndicator from "./LoadingIndicator";
import NoAccess from "./NoAccess";
import HappinessDoor from "./HappinessDoor";

export default class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      connected: false,
      error: false,
    }
  }

  backendUrl() {
    const searchWithAuthTokens = window.location.search;
    return 'ws://localhost:8081' + searchWithAuthTokens
  }

  registerHandlers(socket) {
    socket.on(MessageType.fromBackend.happinessDoorData, (data) => this.happinessDoorData(data));
  }

  componentDidMount() {
    let socket = new Socket(this.backendUrl());
    this.setState({socket})

    socket.on('connect', () => this.onConnect());
    socket.on('disconnect', () => this.onDisconnect());
    socket.on('error', () => this.onError());
    this.registerHandlers(socket);
  }

  onConnect() {
    this.setState({connected: true});
  }

  onDisconnect() {
    this.setState({connected: false});
  }

  onError() {
    this.setState({error: true});
  }

  // helloFromClient() {
  //   if (this.state.connected) {
  //     this.state.socket.emit(MessageType.toBackend.helloFromClient, 'hello server!');
  //   }
  // }

  happinessDoorData(data) {
    this.setState({happinessDoor: data})
  }

  renderLoadingIndicator() {
    return (
      <LoadingIndicator/>
    )
  }

  renderAccessDenied() {
    return (
      <NoAccess/>
    )
  }

  render() {
    if (this.state.error) {
      return this.renderAccessDenied();
    }
    if (!this.state.connected || !this.state.happinessDoor) {
      return this.renderLoadingIndicator();
    }
    return (
      <HappinessDoor meetingName={this.state.happinessDoor.Name}/>
    )
  }
}